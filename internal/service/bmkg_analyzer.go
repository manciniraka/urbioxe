package service

import (
	"time"

	"github.com/manciniraka/urbioxe/internal/entity"
)

type BMKGPeriod struct {
	Weather     string
	Temperature float64
	Humidity    int

	StartTime time.Time
	EndTime   time.Time
}

type BMKGAnalysis struct {
	Current  BMKGPeriod
	Upcoming *BMKGPeriod
}

func AnalyzeForecast(forecasts []entity.WeatherForecast) *BMKGAnalysis {
	if len(forecasts) == 0 {
		return nil
	}

	now := time.Now()

	currentIndex := 0

	for i, forecast := range forecasts {
		forecastTime := forecast.ForecastTime

		if forecastTime.After(now) {
			if i > 0 {
				currentIndex = i - 1
			}

			break
		}

		currentIndex = i
	}

	current := BMKGPeriod{
		Weather: forecasts[currentIndex].Weather,
		Temperature: forecasts[currentIndex].Temperature,
		Humidity: forecasts[currentIndex].Humidity,
		StartTime: forecasts[currentIndex].ForecastTime,
		EndTime: findWeatherEndTime(
			forecasts,
			currentIndex,
		),
	}

	var upcoming *BMKGPeriod

	for i := currentIndex + 1; i < len(forecasts); i++ {
		if forecasts[i].Weather == current.Weather {
			continue
		}

		next := BMKGPeriod{
			Weather: forecasts[i].Weather,
			Temperature: forecasts[i].Temperature,
			Humidity: forecasts[i].Humidity,
			StartTime: forecasts[i].ForecastTime,
			EndTime: findWeatherEndTime(
				forecasts,
				i,
			),
		}

		upcoming = &next

		break
	}

	return &BMKGAnalysis{
		Current:  current,
		Upcoming: upcoming,
	}
}

func findWeatherEndTime(
	forecasts []entity.WeatherForecast,
	startIndex int,
) time.Time {
	currentWeather := forecasts[startIndex].Weather

	for i := startIndex + 1; i < len(forecasts); i++ {
		if forecasts[i].Weather != currentWeather {
			return forecasts[i].ForecastTime
		}
	}

	return forecasts[len(forecasts)-1].ForecastTime
}