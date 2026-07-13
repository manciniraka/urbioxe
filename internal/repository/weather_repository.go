package repository

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WeatherRepository interface {
	GetAllWeather() ([]entity.WeatherCache, error)
	GetWeatherByDistrictID(districtID uint) (*entity.WeatherCache, error)
	Upsert(weather *entity.WeatherCache) error
	GetAllDistricts() ([]entity.District, error)
}

type weatherRepository struct {
	db *gorm.DB
}

func NewWeatherRepository(
	db *gorm.DB,
) WeatherRepository {
	return &weatherRepository{
		db: db,
	}
}

func (wr *weatherRepository) GetAllWeather() ([]entity.WeatherCache, error) {
	var weather []entity.WeatherCache

	err := wr.db.
		// Joins("District"). // TODO: use this after District Module has been merged
		Preload("District").
		// Order("districts.name ASC").
		Find(&weather).Error

	if err != nil {
		return nil, err
	}

	return weather, nil
}

func (wr *weatherRepository) GetWeatherByDistrictID(districtID uint) (*entity.WeatherCache, error) {
	var weather entity.WeatherCache

	err := wr.db.
		Preload("District").
		Where(
			"district_id = ?",
			districtID,
		).
		First(&weather).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrDistrictNotFound
		}

		return nil, err
	}

	return &weather, nil
}

func (wr *weatherRepository) Upsert(weather *entity.WeatherCache) error {
	return wr.db.
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{
						Name: "district_id",
					},
				},
				DoUpdates: clause.AssignmentColumns(
					[]string{
						"temperature",
						"humidity",
						"weather",
						"air_quality",
						"updated_at",
					},
				),
			},
		).
		Create(weather).Error
}

func (wr *weatherRepository) GetAllDistricts() ([]entity.District, error) {
	var districts []entity.District

	err := wr.db.
		Where(
			"is_active = ?",
			true,
		).
		Order("name ASC").
		Find(&districts).Error

	if err != nil {
		return nil, err
	}

	return districts, nil
}
