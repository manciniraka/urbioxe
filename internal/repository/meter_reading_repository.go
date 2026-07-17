package repository

import (
	"time"

	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type MeterReadingRepository interface {
	CreateMeterReading(meterReading *entity.MeterReading) error
	GetMyMeterReadings(userID uint) ([]entity.MeterReading, error)
	GetAllMeterReadings(month *int, year *int) ([]entity.MeterReading, error)
}

type meterReadingRepository struct {
	db *gorm.DB
}

func NewMeterReadingRepository(
	db *gorm.DB,
) MeterReadingRepository {
	return &meterReadingRepository{
		db: db,
	}
}

func (mr *meterReadingRepository) CreateMeterReading(meterReading *entity.MeterReading) error {
	return mr.db.Create(meterReading).Error
}

func (mr *meterReadingRepository) GetMyMeterReadings(userID uint) ([]entity.MeterReading, error) {
	var meterReadings []entity.MeterReading

	err := mr.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&meterReadings).
		Error

	if err != nil {
		return nil, err
	}

	return meterReadings, nil
}

func (mr *meterReadingRepository) GetAllMeterReadings(month *int, year *int) ([]entity.MeterReading, error) {
	var meterReadings []entity.MeterReading

	db := mr.db.Preload("User")

	now := time.Now()

	switch {
		case month != nil && year != nil:
			db = db.
				Where("EXTRACT(MONTH FROM created_at) = ?", *month).
				Where("EXTRACT(YEAR FROM created_at) = ?", *year)

		case year != nil:
			db = db.
				Where("EXTRACT(YEAR FROM created_at) = ?", *year)

		default:
			db = db.
				Where("EXTRACT(MONTH FROM created_at) = ?", int(now.Month())).
				Where("EXTRACT(YEAR FROM created_at) = ?", now.Year())
	}

	err := db.
		Order(`
			CASE status
				WHEN 'Pending' THEN 1
				WHEN 'Approved' THEN 2
				WHEN 'Rejected' THEN 3
			END,
			created_at DESC
		`).
		Find(&meterReadings).
		Error

	if err != nil {
		return nil, err
	}

	return meterReadings, nil
}