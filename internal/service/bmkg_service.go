package service

import (
	"log/slog"
	"time"

	"github.com/manciniraka/urbioxe/external/bmkg"
	"github.com/manciniraka/urbioxe/internal/constant"
	"github.com/manciniraka/urbioxe/internal/dto"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/logger"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type BMKGService interface {
	SyncForecasts() (*dto.SyncForecastResponse, error)
	GetAllWeather() (*dto.WeatherListResult, error)
	GetWeatherDetailsByDistrictID(districtID uint) (*dto.WeatherSummaryResponse, error)
	GetMyWeather(userID uint) (*dto.WeatherSummaryResponse, error)
}

type bmkgService struct {
	db           *gorm.DB
	bmkgClient   *bmkg.Client
	bmkgRepo     repository.BMKGRepository
	districtRepo repository.DistrictRepository
	userRepo     repository.UserRepository
}

func NewBMKGService(
	db *gorm.DB,
	bmkgClient *bmkg.Client,
	bmkgRepo repository.BMKGRepository,
	districtRepo repository.DistrictRepository,
	userRepo repository.UserRepository,
) BMKGService {
	return &bmkgService{
		db: db,
		bmkgClient: bmkgClient,
		bmkgRepo: bmkgRepo,
		districtRepo: districtRepo,
		userRepo: userRepo,
	}
}

func (bs *bmkgService) SyncForecasts() (*dto.SyncForecastResponse, error) {
	districts, err := bs.districtRepo.GetAll()
	if err != nil {
		return nil, err
	}
	

	successCount := 0
	failedCount := 0

	for _, district := range districts {
		logger.Log.Info(
			"district loaded",
			"name", district.Name,
			"adm4", district.BMKGADM4Code,
		)

		forecasts, err := bs.bmkgClient.GetForecast(district.BMKGADM4Code)
		if err != nil {
			failedCount++

			logger.Log.Error(
				"failed to fetch BMKG forecast",
				"tag", constant.LogTagBMKG,
				"district", district.Name,
				"adm4", district.BMKGADM4Code,
				"error", err,
			)

			continue
		}

		var weatherForecasts []entity.WeatherForecast

		for _, forecast := range forecasts {

			weatherForecasts = append(
				weatherForecasts,
				entity.WeatherForecast{
					DistrictID: district.ID,
					ForecastTime: forecast.ForecastTime,
					Temperature: forecast.Temperature,
					Humidity: forecast.Humidity,
					Weather: forecast.Weather,
				},
			)
		}

		err = bs.db.Transaction(
			func(tx *gorm.DB) error {

				if err := bs.bmkgRepo.DeleteForecastsByDistrictTx(
					tx,
					district.ID,
				); err != nil {
					return err
				}

				if err := bs.bmkgRepo.SaveForecastsTx(
					tx,
					weatherForecasts,
				); err != nil {
					return err
				}

				return nil
			},
		)

		if err != nil {
			failedCount++

			logger.Log.Error(
				"failed to save weather forecast",
				"tag", constant.LogTagBMKG,
				"district", district.Name,
				"adm4", district.BMKGADM4Code,
				"error", err,
			)

			continue
		}

		successCount++
	}

	logger.Log.Info(
		"BMKG forecast synchronization completed",
		slog.String("tag", constant.LogTagBMKG),
		slog.Int("success", successCount),
		slog.Int("failed", failedCount),
	)

	return &dto.SyncForecastResponse{
		Success: successCount,
		Failed: failedCount,
		Total: len(districts),
	}, nil
}

func (bs *bmkgService) GetAllWeather() (*dto.WeatherListResult, error) {
	districts, err := bs.districtRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var weatherList []dto.WeatherListResponse

	var latestForecast string

	for _, district := range districts {
		forecasts, err := bs.bmkgRepo.GetForecastByDistrictID(district.ID)
		if err != nil {
			continue
		}

		analysis := AnalyzeForecast(forecasts)
		if analysis == nil {
			continue
		}

		if latestForecast == "" {
			latestForecast = analysis.Current.StartTime.Format(
				"2006-01-02 15:04 WIB",
			)
		}

		weatherList = append(
			weatherList,
			dto.WeatherListResponse{
				District: district.Name,
				Weather: analysis.Current.Weather,
				Temperature: analysis.Current.Temperature,
				Humidity: analysis.Current.Humidity,
			},
		)
	}

	return &dto.WeatherListResult{
		Metadata: dto.WeatherListMetadata{
			ForecastTime: latestForecast,
		},
		Data: weatherList,
	}, nil
}

func (bs *bmkgService) GetWeatherDetailsByDistrictID(districtID uint) (*dto.WeatherSummaryResponse, error) {
	district, err := bs.districtRepo.GetByID(districtID)
	if err != nil {
		return nil, err
	}

	forecasts, err := bs.bmkgRepo.GetForecastByDistrictID(districtID)
	if err != nil {
		return nil, err
	}

	analysis := AnalyzeForecast(forecasts)
	if analysis == nil {
		return nil, errs.ErrWeatherNotFound
	}

	now := time.Now()

	var timeline []dto.ForecastTimelineResponse

	forecastDate := ""

	for _, item := range forecasts {
		forecastTime := item.ForecastTime

		if forecastTime.Before(now) {
			continue
		}

		if !sameDay(forecastTime, now) {
			continue
		}

		timeline = append(
			timeline,
			dto.ForecastTimelineResponse{
				Time: forecastTime.Format("15:04"),
				Weather: item.Weather,
				Temperature: item.Temperature,
				Humidity: item.Humidity,
			},
		)
	}

	if len(timeline) == 0 {
		nextDay := now.AddDate(
			0,
			0,
			1,
		)

		for _, item := range forecasts {
			forecastTime := item.ForecastTime

			if !sameDay(
				forecastTime,
				nextDay,
			) {
				continue
			}

			timeline = append(
				timeline,
				dto.ForecastTimelineResponse{
					Time: forecastTime.Format("15:04"),
					Weather: item.Weather,
					Temperature: item.Temperature,
					Humidity: item.Humidity,
				},
			)

			if forecastDate == "" {	
				forecastDate = forecastTime.Format("Monday, 02 January 2006")
			}
		}
	}

	return &dto.WeatherSummaryResponse{
		District: district.Name,
		ForecastDate: forecastDate,
		Current: dto.CurrentWeatherResponse{
			Weather: analysis.Current.Weather,
			Temperature: analysis.Current.Temperature,
			Humidity: analysis.Current.Humidity,
			ValidUntil: analysis.Current.EndTime.
				Format("15:04"),
		},

		Forecast: timeline,
	}, nil
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()
}

func (bs *bmkgService) GetMyWeather(userID uint) (*dto.WeatherSummaryResponse, error) {
	user, err := bs.userRepo.GetByID(
		userID,
	)
	if err != nil {
		return nil, err
	}

	if user.HomeDistrictID == nil {
		return nil, errs.ErrDistrictNotFound
	}

	return bs.GetWeatherDetailsByDistrictID(*user.HomeDistrictID)
}