package service

import (
	"errors"
	"mime/multipart"

	"github.com/manciniraka/urbioxe/external/cloudinary"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type ReportService interface {
	CreateReport(reportInput entity.Report, files []*multipart.FileHeader) (*entity.Report, error)
	GetAllReports(param GetReportsParam) (*ReportListResponse, error)
	GetReportByID(reportID int64, userID int64, role string) (*entity.Report, error)
	UpdateReport(reportID int64, userID int64, input UpdateReportInput) (*entity.Report, error)
	AssignReport(reportID int64, adminUserID int64, role string, input AssignReportInput) error
	StartReport(reportID int64, officerUserID int64, role string, notes string) error
	ResolveReport(reportID int64, officerUserID int64, role string, notes string, files []*multipart.FileHeader) error
}

type reportService struct {
	repo   repository.ReportRepository
	cldSvc cloudinary.CloudinaryService
}

func NewReportService(
	reportRepo repository.ReportRepository,
	cldSvc cloudinary.CloudinaryService,
) ReportService {
	return &reportService{
		repo:   reportRepo,
		cldSvc: cldSvc,
	}
}

type GetReportsParam struct {
	UserID     int64
	Role       string
	DistrictID int64
	Status     string
	Page       int
	Limit      int
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
	CategoryID         int64    `json:"category_id"`
	IncidentDistrictID int64    `json:"incident_district_id"`
	Latitude           *float64 `json:"latitude"`
	Longitude          *float64 `json:"longitude"`
}

type AssignReportInput struct {
	StaffID int64  `json:"staff_id"`
	Notes   string `json:"notes"`
}

type StartReportInput struct {
	Notes string `json:"notes"`
}

func (rs *reportService) CreateReport(reportInput entity.Report, files []*multipart.FileHeader) (*entity.Report, error) {
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

	return &reportInput, nil
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

func (rs *reportService) GetReportByID(reportID int64, userID int64, role string) (*entity.Report, error) {
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

	return report, nil
}

func (rs *reportService) UpdateReport(reportID int64, userID int64, input UpdateReportInput) (*entity.Report, error) {
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

	return existingReport, nil
}

func (rs *reportService) AssignReport(reportID int64, adminUserID int64, role string, input AssignReportInput) error {
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

	return rs.repo.AssignStaff(reportID, input.StaffID, adminUserID, notes)
}

func (rs *reportService) StartReport(reportID int64, officerUserID int64, role string, notes string) error {
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

	if report.Status != entity.StatusAssigned {
		if report.Status == entity.StatusInProgress {
			return errs.ErrReportAlreadyInProcess
		}
		return errs.ErrReportNotAssigned
	}

	if role == "officer" {
		if report.AssignedStaffID == nil {
			return errs.ErrReportNotAssigned
		}
	}

	if notes == "" {
		notes = "Officer start handle the report"
	}

	return rs.repo.UpdateStatusWithHistory(
		reportID,
		entity.StatusInProgress,
		officerUserID,
		notes,
		false,
	)
}

func (rs *reportService) ResolveReport(reportID int64, officerUserID int64, role string, notes string, files []*multipart.FileHeader) error {
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

	return rs.repo.ResolveReport(reportID, officerUserID, notes, resolutionAttachments)
}
