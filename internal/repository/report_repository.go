package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(report *entity.Report) error
	FindAll(filter ReportFilter) ([]entity.Report, int64, error)
	UpdateStatusWithHistory(reportID int64, newStatus entity.ReportStatus, actorID int64, notes string, isInternal bool) error
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}

type ReportFilter struct {
	UserID     int64
	Role       string
	DistrictID int64
	Status     string
	Limit      int
	Offset     int
}

func (rr *reportRepository) Create(report *entity.Report) error {
	return rr.db.Create(report).Error
}

func (rr *reportRepository) UpdateStatusWithHistory(reportID int64, newStatus entity.ReportStatus, actorID int64, notes string, isInternal bool) error {
	return rr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Report{}).Where("id = ?", reportID).Update("status", newStatus).Error; err != nil {
			return err
		}

		history := entity.ReportHistory{
			ReportID:   reportID,
			Status:     newStatus,
			Notes:      notes,
			IsInternal: isInternal,
			ActorID:    &actorID,
		}

		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}

func (rr *reportRepository) FindAll(filter ReportFilter) ([]entity.Report, int64, error) {
	var reports []entity.Report
	var totalData int64

	query := rr.db.Model(&entity.Report{}).
		Preload("Category").
		Preload("District").
		Preload("User").
		Preload("Attachments").
		Preload("Histories")

	if filter.Role == "citizen" {
		query = query.Where("reports.user_id = ?", filter.UserID)
	}

	if filter.Role == "officer" {
		query = query.Where(
			"reports.assigned_staff_id IN (SELECT id FROM staff_profiles WHERE user_id = ?)",
			filter.UserID,
		)
	}

	if filter.Role == "department_admin" {
		query = query.Joins("JOIN categories ON categories.id = reports.category_id").
			Where("categories.department_id = (SELECT department_id FROM staff_profiles WHERE user_id = ?)", filter.UserID)
	}

	if filter.DistrictID != 0 {
		query = query.Where("reports.incident_district_id = ?", filter.DistrictID)
	}

	if filter.Status != "" {
		query = query.Where("reports.status = ?", filter.Status)
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(filter.Limit).Offset(filter.Offset).Order("reports.created_at DESC").Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, totalData, nil
}
