package service

import (
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"time"

	"github.com/manciniraka/urbioxe/external/cloudinary"
	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/constant"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/logger"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type ReportService interface {
	CreateReport(reportInput entity.Report, files []*multipart.FileHeader) (*ReportResponse, error)
	GetAllReports(param GetReportsParam) (*ReportListResponse, error)
	GetReportByID(reportID uint, userID uint, role string) (*ReportResponse, error)
	UpdateReport(reportID uint, userID uint, input UpdateReportInput) (*ReportResponse, error)
	AssignReport(reportID uint, adminUserID uint, role string, input AssignReportInput) error
	StartReport(reportID uint, officerUserID uint, role string, notes string) error
	ResolveReport(reportID uint, officerUserID uint, role string, notes string, files []*multipart.FileHeader) error
	RejectReport(reportID uint, adminUserID uint, role string, input UpdateStatusReportInput) error
	VerifyReport(reportID uint, adminUserID uint, role string, input UpdateStatusReportInput) error
	UpdatePriority(reportID uint, adminUserID uint, role string, input UpdatePriorityInput) error
	ReassignReport(reportID uint, adminUserID uint, role string, input ReassignReportInput) error
}

type reportService struct {
	repo      repository.ReportRepository
	staffRepo repository.StaffRepository
	cldSvc    cloudinary.CloudinaryService
	mailer    *mailjet.Client
}

func NewReportService(
	reportRepo repository.ReportRepository,
	staffRepo repository.StaffRepository,
	cldSvc cloudinary.CloudinaryService,
	mailer *mailjet.Client,
) ReportService {
	return &reportService{
		repo:      reportRepo,
		staffRepo: staffRepo,
		cldSvc:    cldSvc,
		mailer:    mailer,
	}
}

type GetReportsParam struct {
	UserID     uint
	Role       string
	DistrictID uint
	Status     string
	Page       int
	Limit      int
}

type UserReportResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReportResponse struct {
	ID              uint                      `json:"id"`
	Title           string                    `json:"title"`
	Description     string                    `json:"description"`
	AddressLandmark string                    `json:"address_landmark"`
	Status          entity.ReportStatus       `json:"status"`
	Priority        entity.ReportPriority     `json:"priority"`
	CreatedAt       time.Time                 `json:"created_at"`
	User            *UserReportResponse       `json:"user,omitempty"`
	Category        *entity.Category          `json:"category,omitempty"`
	Attachments     []entity.ReportAttachment `json:"attachments,omitempty"`
	Histories       []entity.ReportHistory    `json:"histories,omitempty"`
}

type ReportListResponse struct {
	Data      []entity.Report `json:"data"`
	TotalData int64           `json:"total_data"`
	Page      int             `json:"page"`
	Limit     int             `json:"limit"`
	TotalPage int             `json:"total_page"`
}

type UpdateReportInput struct {
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	AddressLandmark    string   `json:"address_landmark"`
	CategoryID         uint     `json:"category_id"`
	IncidentDistrictID uint     `json:"incident_district_id"`
	Latitude           *float64 `json:"latitude"`
	Longitude          *float64 `json:"longitude"`
}

type AssignReportInput struct {
	StaffID uint   `json:"staff_id"`
	Notes   string `json:"notes"`
}

type UpdateStatusReportInput struct {
	Notes string `json:"notes"`
}

type UpdatePriorityInput struct {
	Priority entity.ReportPriority `json:"priority"`
	Notes    string                `json:"notes"`
}

type ReassignReportInput struct {
	NewStaffID uint   `json:"new_staff_id"`
	Notes      string `json:"notes"`
}

// helper
func ToReportResponse(r *entity.Report) *ReportResponse {
	var userResp *UserReportResponse
	if r.User != nil {
		userResp = &UserReportResponse{
			ID:    r.User.ID,
			Name:  r.User.Name,
			Email: r.User.Email,
		}
	}

	return &ReportResponse{
		ID:              r.ID,
		Title:           r.Title,
		Description:     r.Description,
		AddressLandmark: r.AddressLandmark,
		Status:          r.Status,
		Priority:        r.Priority,
		CreatedAt:       r.CreatedAt,
		User:            userResp,
		Category:        r.Category,
		Attachments:     r.Attachments,
		Histories:       r.Histories,
	}
}

func GenerateReportNumber(districtID uint, categoryID uint) string {
	now := time.Now()
	fullTimestamp := now.Format("20060102150405")
	randomNumber := rand.Intn(900) + 100
	return fmt.Sprintf("REP-%s-D%03d-C%03d-%d", fullTimestamp, districtID, categoryID, randomNumber)
}

// helper send email per status update
func (rs *reportService) sendStatusEmailAsync(report *entity.Report, status string, notes string) {
	if report != nil && report.User != nil && report.User.Email != "" {
		go func(toEmail, toName, title string, id uint, currentStatus, currentNotes string) {
			err := rs.mailer.SendReportStatusEmail(
				toEmail,
				toName,
				title,
				id,
				currentStatus,
				currentNotes,
			)
			if err != nil {
				logger.Log.Error(
					"failed to send update report email",
					"tag", constant.LogTagMailjet,
					"email", toEmail,
					"status", currentStatus,
					"error", err,
				)
			}
		}(
			report.User.Email,
			report.User.Name,
			report.Title,
			report.ID,
			status,
			notes,
		)
	}
}

func (rs *reportService) CreateReport(reportInput entity.Report, files []*multipart.FileHeader) (*ReportResponse, error) {
	var attachments []entity.ReportAttachment

	for _, fileHeader := range files {
		secureURL, err := rs.cldSvc.UploadImage(fileHeader)
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, entity.ReportAttachment{
			FileURL: secureURL,
			Type:    entity.TypeEvidence,
		})
	}

	reportInput.Attachments = attachments
	reportInput.Status = entity.StatusPending
	reportNumber := GenerateReportNumber(reportInput.IncidentDistrictID, reportInput.CategoryID)
	reportInput.ReportNumber = &reportNumber

	reportInput.Histories = []entity.ReportHistory{
		{
			Status:     entity.StatusPending,
			Notes:      "Report already created. Waiting for verification.",
			IsInternal: false,
			ActorID:    &reportInput.UserID,
		},
	}

	err := rs.repo.Create(&reportInput)
	if err != nil {
		return nil, err
	}

	formattedDate := reportInput.CreatedAt.Format("02 January 2006 15:04 MST")

	categoryName := "Umum"
	if reportInput.Category != nil {
		categoryName = reportInput.Category.Name
	}

	if reportInput.User != nil && reportInput.User.Email != "" {
		go func(toEmail, toName string, id uint, title, category, address, dateStr string) {
			errMail := rs.mailer.SendReportCreatedEmail(
				toEmail,
				toName,
				id,
				title,
				category,
				address,
				dateStr,
			)
			if errMail != nil {
				logger.Log.Error(
					"failed to send create report email",
					"tag", constant.LogTagMailjet,
					"email", toEmail,
					"error", errMail,
				)
			}
		}(
			reportInput.User.Email,
			reportInput.User.Name,
			reportInput.ID,
			reportInput.Title,
			categoryName,
			reportInput.AddressLandmark,
			formattedDate,
		)
	} else {
		logger.Log.Error(
			"skip to send create report email",
			"tag", constant.LogTagMailjet,
			"email", reportInput.User.Email,
			"error", err,
		)
	}

	return ToReportResponse(&reportInput), nil
}

func (rs *reportService) GetAllReports(param GetReportsParam) (*ReportListResponse, error) {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	offset := (param.Page - 1) * param.Limit

	repoFilter := repository.ReportFilter{
		UserID:     param.UserID,
		Role:       param.Role,
		DistrictID: param.DistrictID,
		Status:     param.Status,
		Limit:      param.Limit,
		Offset:     offset,
	}

	reports, totalData, err := rs.repo.FindAll(repoFilter)
	if err != nil {
		return nil, err
	}

	totalPage := int(totalData) / param.Limit
	if int(totalData)%param.Limit != 0 {
		totalPage++
	}

	return &ReportListResponse{
		Data:      reports,
		TotalData: totalData,
		Page:      param.Page,
		Limit:     param.Limit,
		TotalPage: totalPage,
	}, nil
}

func (rs *reportService) GetReportByID(reportID uint, userID uint, role string) (*ReportResponse, error) {
	report, err := rs.repo.FindByID(reportID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrReportNotFound
		}
		return nil, err
	}

	if role == "citizen" && report.UserID != userID {
		return nil, errs.ErrReportForbidden
	}

	return ToReportResponse(report), nil
}

func (rs *reportService) UpdateReport(reportID uint, userID uint, input UpdateReportInput) (*ReportResponse, error) {
	existingReport, err := rs.repo.FindByID(reportID, "citizen")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrReportNotFound
		}
		return nil, err
	}

	if existingReport.UserID != userID {
		return nil, errs.ErrReportForbidden
	}

	if existingReport.Status != entity.StatusPending {
		return nil, errs.ErrReportAlreadyInProcess
	}

	if input.Title != "" {
		existingReport.Title = input.Title
	}
	if input.Description != "" {
		existingReport.Description = input.Description
	}
	if input.AddressLandmark != "" {
		existingReport.AddressLandmark = input.AddressLandmark
	}
	if input.CategoryID != 0 {
		existingReport.CategoryID = input.CategoryID
	}
	if input.IncidentDistrictID != 0 {
		existingReport.IncidentDistrictID = input.IncidentDistrictID
	}
	if input.Latitude != nil {
		existingReport.Latitude = input.Latitude
	}
	if input.Longitude != nil {
		existingReport.Longitude = input.Longitude
	}

	err = rs.repo.Update(existingReport)
	if err != nil {
		return nil, err
	}

	return ToReportResponse(existingReport), nil
}

func (rs *reportService) AssignReport(reportID uint, adminUserID uint, role string, input AssignReportInput) error {
	if role != "department_admin" && role != "super_admin" {
		return errs.ErrReportAssignForbidden
	}

	if input.StaffID == 0 {
		return errors.New("staff_id required")
	}

	report, err := rs.repo.FindByID(reportID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrReportNotFound
		}
		return err
	}

	if report.AssignedStaffID != nil || report.Status == entity.StatusAssigned {
		return errs.ErrReportAlreadyAssigned
	}

	if report.Status == entity.StatusInProgress {
		return errs.ErrReportAlreadyInProcess
	}

	if report.Status == entity.StatusResolved || report.Status == entity.StatusRejected {
		return errs.ErrReportAlreadyResolved
	}

	notes := input.Notes
	if notes == "" {
		notes = "Report already assigned to field officer."
	}

	err = rs.repo.AssignStaff(reportID, input.StaffID, adminUserID, notes)
	if err != nil {
		return err
	}

	rs.sendStatusEmailAsync(report, string(entity.StatusAssigned), notes)

	return nil
}

func (rs *reportService) StartReport(reportID uint, officerUserID uint, role string, notes string) error {
	if role != "officer" && role != "super_admin" {
		return errs.ErrReportUpdateForbidden
	}

	report, err := rs.repo.FindByID(reportID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrReportNotFound
		}
		return err
	}

	staff, err := rs.staffRepo.FindStaffByUserID(officerUserID)
	if err != nil {
		return err
	}

	if role == "officer" {
		if report.AssignedStaffID == nil {
			return errs.ErrReportNotAssigned
		}

		if staff.ID != *report.AssignedStaffID {
			return errs.ErrReportUpdateForbidden
		}
	}

	if report.Status != entity.StatusAssigned {
		if report.Status == entity.StatusInProgress {
			return errs.ErrReportAlreadyInProcess
		}
		return errs.ErrReportNotAssigned
	}

	if notes == "" {
		notes = "Officer start handle the report"
	}

	err = rs.repo.UpdateStatusWithHistory(
		reportID,
		entity.StatusInProgress,
		officerUserID,
		notes,
		false,
	)
	if err != nil {
		return err
	}

	rs.sendStatusEmailAsync(report, string(entity.StatusInProgress), notes)

	return nil
}

func (rs *reportService) ResolveReport(reportID uint, officerUserID uint, role string, notes string, files []*multipart.FileHeader) error {
	if role != "officer" && role != "super_admin" {
		return errs.ErrReportUpdateForbidden
	}

	if len(files) == 0 {
		return errors.New("upload minimal 1 image")
	}

	report, err := rs.repo.FindByID(reportID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrReportNotFound
		}
		return err
	}

	staff, err := rs.staffRepo.FindStaffByUserID(officerUserID)
	if err != nil {
		return err
	}

	if role == "officer" {
		if report.AssignedStaffID == nil {
			return errs.ErrReportNotAssigned
		}

		if staff.ID != *report.AssignedStaffID {
			return errs.ErrReportUpdateForbidden
		}
	}

	if report.Status != entity.StatusInProgress {
		return errs.ErrReportShouldInProcess
	}

	var resolutionAttachments []entity.ReportAttachment
	for _, fileHeader := range files {
		secureURL, err := rs.cldSvc.UploadImage(fileHeader)
		if err != nil {
			return errors.New("failed upload image: " + err.Error())
		}

		resolutionAttachments = append(resolutionAttachments, entity.ReportAttachment{
			FileURL: secureURL,
			Type:    entity.TypeResolution,
		})
	}

	if notes == "" {
		notes = "Report mark as solved by officer"
	}

	err = rs.repo.ResolveReport(reportID, officerUserID, notes, resolutionAttachments)
	if err != nil {
		return err
	}

	rs.sendStatusEmailAsync(report, string(entity.StatusResolved), notes)
	return nil
}

func (rs *reportService) RejectReport(reportID uint, adminUserID uint, role string, input UpdateStatusReportInput) error {
	if role != "department_admin" && role != "super_admin" {
		return errs.ErrReportUpdateForbidden
	}

	if input.Notes == "" {
		return errors.New("notes required!")
	}

	report, err := rs.repo.FindByID(reportID, role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrReportNotFound
		}
		return err
	}

	if report.Status == entity.StatusInProgress {
		return errs.ErrReportAlreadyInProcess
	}
	if report.Status == entity.StatusResolved {
		return errs.ErrReportAlreadyResolved
	}
	if report.Status == entity.StatusRejected {
		return errs.ErrReportAlreadyRejected
	}

	err = rs.repo.UpdateStatusWithHistory(
		reportID,
		entity.StatusRejected,
		adminUserID,
		input.Notes,
		false,
	)

	if err != nil {
		return err
	}

	rs.sendStatusEmailAsync(report, string(entity.StatusRejected), input.Notes)
	return nil
}

func (rs *reportService) VerifyReport(reportID uint, adminUserID uint, role string, input UpdateStatusReportInput) error {
	if role != "admin" && role != "super_admin" {
		return errs.ErrReportUpdateForbidden
	}

	report, err := rs.repo.FindByID(reportID, role)
	if err != nil {
		return err
	}

	if report.Status != entity.StatusPending {
		return errs.ErrReportShouldPending
	}

	notes := input.Notes
	if notes == "" {
		notes = "Report already verified by admin and it's valid."
	}

	err = rs.repo.UpdateStatusWithHistory(reportID, entity.StatusVerified, adminUserID, notes, false)
	if err != nil {
		return err
	}

	rs.sendStatusEmailAsync(report, string(entity.StatusVerified), notes)

	return nil
}

func (rs *reportService) UpdatePriority(reportID uint, adminUserID uint, role string, input UpdatePriorityInput) error {
	return nil
}

func (rs *reportService) ReassignReport(reportID uint, adminUserID uint, role string, input ReassignReportInput) error {
	return nil
}
