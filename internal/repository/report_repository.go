package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(report *entity.Report) error
	FindAll(filter ReportFilter) ([]entity.Report, int64, error)
	UpdateStatusWithHistory(reportID uint, newStatus entity.ReportStatus, actorID uint, notes string, isInternal bool) error
	FindByID(id uint, role string) (*entity.Report, error)
	Update(report *entity.Report) error
	AssignStaff(reportID uint, staffID uint, actorID uint, notes string) error
	ResolveReport(reportID uint, actorID uint, notes string, attachments []entity.ReportAttachment) error
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
	UserID     uint
	Role       string
	DistrictID uint
	Status     string
	Limit      int
	Offset     int
}

func (rr *reportRepository) Create(report *entity.Report) error {
	if err := rr.db.Create(report).Error; err != nil {
		return err
	}

	return rr.db.Preload("User").
		Preload("Category").
		Preload("Attachments").
		First(report, report.ID).Error
}

func (rr *reportRepository) UpdateStatusWithHistory(reportID uint, newStatus entity.ReportStatus, actorID uint, notes string, isInternal bool) error {
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

func (rr *reportRepository) FindByID(id uint, role string) (*entity.Report, error) {
	var report entity.Report

	query := rr.db.Model(&entity.Report{}).
		Preload("Category").
		Preload("District").
		Preload("User").
		Preload("Attachments")

	if role == "citizen" {
		query = query.Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_internal = ?", false).Order("report_histories.created_at ASC")
		})
	} else {
		query = query.Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Order("report_histories.created_at ASC")
		})
	}

	err := query.First(&report, id).Error
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (rr *reportRepository) Update(report *entity.Report) error {
	return rr.db.Save(report).Error
}

func (rr *reportRepository) AssignStaff(reportID uint, staffID uint, actorID uint, notes string) error {
	return rr.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"assigned_staff_id": staffID,
			"status":            entity.StatusAssigned,
		}
		if err := tx.Model(&entity.Report{}).Where("id = ?", reportID).Updates(updates).Error; err != nil {
			return err
		}

		history := entity.ReportHistory{
			ReportID:   reportID,
			Status:     entity.StatusAssigned,
			Notes:      notes,
			IsInternal: false,
			ActorID:    &actorID,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}

func (rr *reportRepository) ResolveReport(reportID uint, actorID uint, notes string, attachments []entity.ReportAttachment) error {
	return rr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Report{}).Where("id = ?", reportID).Update("status", entity.StatusResolved).Error; err != nil {
			return err
		}

		for i := range attachments {
			attachments[i].ReportID = reportID
		}
		if len(attachments) > 0 {
			if err := tx.Create(&attachments).Error; err != nil {
				return err
			}
		}

		history := entity.ReportHistory{
			ReportID:   reportID,
			Status:     entity.StatusResolved,
			Notes:      notes,
			IsInternal: false,
			ActorID:    &actorID,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}
