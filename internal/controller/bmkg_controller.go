package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type BMKGController interface {
	SyncForecasts(c echo.Context) error
	GetAllWeather(c echo.Context) error
	GetWeatherDetailsByDistrictID(c echo.Context) error
	GetMyWeather(c echo.Context) error
}

type bmkgController struct {
	bmkgService service.BMKGService
}

func NewBMKGController(
	bmkgService service.BMKGService,
) BMKGController {

	return &bmkgController{
		bmkgService: bmkgService,
	}
}

func (bc *bmkgController) SyncForecasts(c echo.Context) error {
	response, err := bc.bmkgService.SyncForecasts()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"weather forecast synchronized successfully",
		response,
	)
}

func (bc *bmkgController) GetAllWeather(c echo.Context) error {
	weatherList, err := bc.bmkgService.GetAllWeather()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"message": "weather list retrieved successfully",
			"metadata": weatherList.Metadata,
			"data": weatherList.Data,
		},
	)
}

func (bc *bmkgController) GetWeatherDetailsByDistrictID(c echo.Context) error {
	districtID, err := strconv.Atoi(c.Param("district_id"))
	if err != nil {
		return helper.HandleError(
			c,
			errs.ErrBadRequest,
		)
	}

	response, err := bc.bmkgService.GetWeatherDetailsByDistrictID(uint(districtID))
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"weather summary retrieved successfully",
		response,
	)
}

func (bc *bmkgController) GetMyWeather(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	response, err := bc.bmkgService.GetMyWeather(userID)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"weather summary retrieved successfully",
		response,
	)
}