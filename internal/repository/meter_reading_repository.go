package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type MeterReadingRepository interface {
	CreateMeterReading(meterReading *entity.MeterReading) error
	GetMyMeterReadings(userID uint) ([]entity.MeterReading, error)
	GetAllMeterReadings() ([]entity.MeterReading, error)
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

func (mr *meterReadingRepository) GetAllMeterReadings() ([]entity.MeterReading, error) {
	var meterReadings []entity.MeterReading

	err := mr.db.
		Preload("User").
		Order("created_at DESC").
		Find(&meterReadings).
		Error

	if err != nil {
		return nil, err
	}

	return meterReadings, nil
}