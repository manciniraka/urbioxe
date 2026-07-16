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

// SyncForecasts godoc
//
//	@Summary		Sync weather forecasts
//	@Description	Synchronize weather forecasts from BMKG
//	@Tags			BMKG
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	helper.Response
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/weather/sync [post]
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

// GetAllWeather godoc
//
//	@Summary		Get all weather
//	@Description	Retrieve weather information for all districts
//	@Tags			BMKG
//	@Produce		json
//	@Success		200	{object}	helper.Response
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/weather [get]
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
			"message":  "weather list retrieved successfully",
			"metadata": weatherList.Metadata,
			"data":     weatherList.Data,
		},
	)
}

// GetWeatherDetailsByDistrictID godoc
//
//	@Summary		Get weather by district
//	@Description	Retrieve weather details by district ID
//	@Tags			BMKG
//	@Produce		json
//	@Param			district_id	path		int	true	"District ID"
//	@Success		200			{object}	helper.Response
//	@Failure		400			{object}	helper.ErrorResponse
//	@Failure		404			{object}	helper.ErrorResponse
//	@Router			/weather/{district_id} [get]
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

// GetMyWeather godoc
//
//	@Summary		Get my weather
//	@Description	Retrieve weather based on authenticated user's district
//	@Tags			BMKG
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	helper.Response
//	@Failure		401	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/weather/me [get]
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
