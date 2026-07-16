package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"

	"gorm.io/gorm"
)

type BMKGRepository interface {
	SaveForecastsTx(tx *gorm.DB, forecasts []entity.WeatherForecast) error
	DeleteForecastsByDistrictTx(tx *gorm.DB, districtID uint) error
	GetForecastByDistrictID(districtID uint) ([]entity.WeatherForecast, error)
}

type bmkgRepository struct {
	db *gorm.DB
}

func NewBMKGRepository(
	db *gorm.DB,
) BMKGRepository {
	return &bmkgRepository{
		db: db,
	}
}

func (br *bmkgRepository) SaveForecastsTx(tx *gorm.DB, forecasts []entity.WeatherForecast) error {
	if len(forecasts) == 0 {
		return nil
	}

	return tx.Create(&forecasts).Error
}

func (br *bmkgRepository) DeleteForecastsByDistrictTx(tx *gorm.DB, districtID uint) error {
	return tx.
		Where(
			"district_id = ?",
			districtID,
		).
		Delete(
			&entity.WeatherForecast{},
		).
		Error
}

func (br *bmkgRepository) GetForecastByDistrictID(districtID uint) ([]entity.WeatherForecast, error) {
	var forecasts []entity.WeatherForecast

	err := br.db.
		Where("district_id = ?", districtID).
		Order("forecast_time ASC").
		Find(&forecasts).
		Error

	if err != nil {
		return nil, err
	}

	return forecasts, nil
}
