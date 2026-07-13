package service

import (
	"time"

	"github.com/manciniraka/urbioxe/external/openweather"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/logger"
	"github.com/manciniraka/urbioxe/internal/repository"
)

type WeatherService interface {
	GetAllWeather() ([]WeatherSummary, error)
	GetWeatherByDistrictID(districtID uint) (*entity.WeatherCache, error)
	GetMyWeather(userID uint) (*entity.WeatherCache, error)
	SyncWeather() (*SyncResult, error)
}

type weatherService struct {
	weatherRepo repository.WeatherRepository
	userRepo    repository.UserRepository

	openWeather *openweather.Client
}

func NewWeatherService(
	weatherRepo repository.WeatherRepository,
	userRepo repository.UserRepository,
	openWeather *openweather.Client,
) WeatherService {

	return &weatherService{
		weatherRepo: weatherRepo,
		userRepo:    userRepo,
		openWeather: openWeather,
	}
}

type WeatherSummary struct {
	DistrictID   uint   `json:"district_id"`
	DistrictName string `json:"district_name"`
	Weather      string `json:"weather"`

	UpdatedAt time.Time `json:"-"`
}

type SyncResult struct {
	TotalDistrict int `json:"total_district"`
	SuccessCount  int `json:"success_count"`
	FailedCount   int `json:"failed_count"`
}

func (ws *weatherService) SyncWeather() (*SyncResult, error) {
	districts, err := ws.weatherRepo.GetAllDistricts()
	if err != nil {
		return nil, err
	}

	syncResult := &SyncResult{
		TotalDistrict: len(districts),
	}

	for _, district := range districts {

		currentWeather, err := ws.openWeather.GetCurrentWeather(
			district.Name,
		)
		if err != nil {

			logger.Log.Error(
				"failed to fetch weather data",
				"district_id",
				district.ID,
				"district",
				district.Name,
				"error",
				err,
			)

			syncResult.FailedCount++

			continue
		}

		weatherCache := entity.WeatherCache{
			DistrictID:  district.ID,
			Temperature: currentWeather.Temperature,
			Humidity:    currentWeather.Humidity,
			Weather:     currentWeather.Weather,
			AirQuality:  currentWeather.AirQuality,
		}

		if err := ws.weatherRepo.Upsert(
			&weatherCache,
		); err != nil {

			logger.Log.Error(
				"failed to save weather cache",
				"district",
				district.Name,
				"error",
				err,
			)

			syncResult.FailedCount++

			continue
		}

		syncResult.SuccessCount++
	}

	return syncResult, nil
}

func (ws *weatherService) GetAllWeather() ([]WeatherSummary, error) {
	weather, err := ws.weatherRepo.GetAllWeather()
	if err != nil {
		return nil, err
	}

	summaries := make(
		[]WeatherSummary,
		0,
		len(weather),
	)

	for _, w := range weather {

		summaries = append(
			summaries,
			WeatherSummary{
				DistrictID:   w.DistrictID,
				DistrictName: w.District.Name,
				Weather:      w.Weather,
				UpdatedAt:    w.UpdatedAt,
			},
		)
	}

	return summaries, nil
}

func (ws *weatherService) GetWeatherByDistrictID(districtID uint) (*entity.WeatherCache, error) {
	return ws.weatherRepo.GetWeatherByDistrictID(
		districtID,
	)
}

func (ws *weatherService) GetMyWeather(userID uint) (*entity.WeatherCache, error) {
	user, err := ws.userRepo.GetByID(
		userID,
	)
	if err != nil {
		return nil, err
	}

	logger.Log.Info(
		"debug user",
		"user_id", user.ID,
		"home_district_nil", user.HomeDistrictID == nil,
	)

	if user.HomeDistrictID == nil {
		logger.Log.Info("return ErrHomeDistrictNotSet")
		return nil, errs.ErrHomeDistrictNotSet
	}

	return ws.weatherRepo.GetWeatherByDistrictID(
		*user.HomeDistrictID,
	)
}
