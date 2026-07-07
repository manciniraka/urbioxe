package service

import (
	"mime/multipart"

	"github.com/manciniraka/urbioxe/external/cloudinary"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/repository"
)

type ReportService interface {
	CreateReport(reportInput entity.Report, files []*multipart.FileHeader) (*entity.Report, error)
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

	err := rs.repo.Create(&reportInput)
	if err != nil {
		return nil, err
	}

	return &reportInput, nil
}
