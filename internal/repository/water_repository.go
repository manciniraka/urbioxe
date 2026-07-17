package repository

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"gorm.io/gorm"
)

type WaterRepository interface {
	CreateWaterStatus(waterStatus *entity.WaterStatus) error
	GetLatestAllWaterStatus() ([]entity.WaterStatus, error)
	GetLatestWaterStatusByDistrictID(districtID uint) (*entity.WaterStatus, error)
	GetWaterStatusHistories() ([]entity.WaterStatus, error)
}

type waterRepository struct {
	db *gorm.DB
}

func NewWaterRepository(
	db *gorm.DB,
) WaterRepository {
	return &waterRepository{
		db: db,
	}
}

func (wr *waterRepository) CreateWaterStatus(waterStatus *entity.WaterStatus) error {
	return wr.db.Create(waterStatus).Error
}

func (wr *waterRepository) GetLatestWaterStatusByDistrictID(districtID uint) (*entity.WaterStatus, error) {
	var waterStatus entity.WaterStatus

	err := wr.db.
		Preload("District").
		Preload("User").
		Where(
			"district_id = ?",
			districtID,
		).
		Order("created_at DESC").
		First(&waterStatus).
		Error

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errs.ErrWaterStatusNotFound
		}

		return nil, err
	}

	return &waterStatus, nil
}


func (wr *waterRepository) GetWaterStatusHistories() ([]entity.WaterStatus, error) {
	var histories []entity.WaterStatus

	err := wr.db.
		Preload("District").
		Preload("User").
		Order("created_at DESC").
		Find(&histories).
		Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (wr *waterRepository) GetLatestAllWaterStatus() ([]entity.WaterStatus, error) {
	var waterStatuses []entity.WaterStatus

	err := wr.db.
		Raw(`
			SELECT DISTINCT ON (district_id) *
			FROM water_statuses
			ORDER BY district_id, created_at DESC
		`).
		Scan(&waterStatuses).
		Error

	if err != nil {
		return nil, err
	}

	for i := range waterStatuses {

		wr.db.
			Model(&waterStatuses[i]).
			Association("District").
			Find(&waterStatuses[i].District)

		wr.db.
			Model(&waterStatuses[i]).
			Association("User").
			Find(&waterStatuses[i].User)
	}

	return waterStatuses, nil
}