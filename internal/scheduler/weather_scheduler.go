package scheduler

import (
	"time"

	"github.com/manciniraka/urbioxe/internal/logger"
	"github.com/manciniraka/urbioxe/internal/service"

	"github.com/robfig/cron/v3"
)

type WeatherScheduler struct {
	weatherService service.WeatherService
}

func NewWeatherScheduler(
	weatherService service.WeatherService,
) *WeatherScheduler {
	return &WeatherScheduler{
		weatherService: weatherService,
	}
}

func (ws *WeatherScheduler) Start() {
	logger.Log.Info(
		"weather scheduler started",
		"cron",
		"0 * * * *",
	)

	c := cron.New(
		cron.WithLocation(
			time.Local,
		),
	)

	// Synchronize weather caches every hour
	_, err := c.AddFunc(
		"0 * * * *",
		func() {
			logger.Log.Info(
				"starting weather synchronization",
			)

			result, err := ws.weatherService.SyncWeather()
			if err != nil {

				logger.Log.Error(
					"weather scheduler failed",
					"error",
					err,
				)

				return
			}

			logger.Log.Info(
				"weather synchronized successfully",
				"total_district",
				result.TotalDistrict,
				"success_count",
				result.SuccessCount,
				"failed_count",
				result.FailedCount,
			)
		},
	)

	if err != nil {

		logger.Log.Error(
			"failed to register weather scheduler",
			"error",
			err,
		)

		return
	}

	c.Start()
}
