package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type WeatherController struct {
	weatherService service.WeatherService
}

func NewWeatherController(
	weatherService service.WeatherService,
) *WeatherController {

	return &WeatherController{
		weatherService: weatherService,
	}
}

func (wc *WeatherController) GetAllWeather(c echo.Context) error {
	weather, err := wc.weatherService.GetAllWeather()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	var updatedAt *time.Time

	if len(weather) > 0 {
		updatedAt = &weather[0].UpdatedAt
	}

	return c.JSON(
		http.StatusOK,
		echo.Map{
			"message":    "weather retrieved successfully",
			"updated_at": updatedAt,
			"data":       weather,
		},
	)
}

func (wc *WeatherController) GetWeatherByDistrictID(c echo.Context) error {
	districtID, err := strconv.Atoi(
		c.Param("district_id"),
	)
	if err != nil {
		return err
	}

	weather, err := wc.weatherService.GetWeatherByDistrictID(uint(districtID))
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"weather retrieved successfully",
		weather,
	)
}

func (wc *WeatherController) GetMyWeather(c echo.Context) error {
	userID := helper.GetUserID(c)

	weather, err := wc.weatherService.GetMyWeather(userID)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"weather retrieved successfully",
		weather,
	)
}

func (wc *WeatherController) SyncWeather(c echo.Context) error {
	result, err := wc.weatherService.SyncWeather()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"weather synchronized successfully",
		result,
	)
}
