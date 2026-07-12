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
