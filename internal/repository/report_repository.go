package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(report *entity.Report) error
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}

func (rr *reportRepository) Create(report *entity.Report) error {
	return rr.db.Create(report).Error
}
