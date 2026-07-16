package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/dto"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type WaterController struct {
	waterService service.WaterService
}

func NewWaterController(
	waterService service.WaterService,
) *WaterController {
	return &WaterController{
		waterService: waterService,
	}
}

func (wc *WaterController) CreateWaterStatus(c echo.Context) error {
	var input dto.CreateWaterStatusInput

	if err := c.Bind(&input); err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	if err := c.Validate(&input); err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	userID := helper.GetUserID(c)

	response, err := wc.waterService.CreateWaterStatus(
		userID,
		input,
	)

	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Created(
		c,
		"water status created successfully",
		response,
	)
}

func (wc *WaterController) GetAllWaterStatus(c echo.Context) error {
	result, err := wc.waterService.GetAllWaterStatus()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"message": "water status retrieved successfully",
			"metadata": result.Metadata,
			"data": result.Data,
		},
	)
}

func (wc *WaterController) GetWaterStatusByDistrictID(c echo.Context) error {
	districtID, err := strconv.Atoi(c.Param("district_id"))
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	result, err := wc.waterService.GetWaterStatusByDistrictID(uint(districtID))
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"water status retrieved successfully",
		result,
	)
}

func (wc *WaterController) GetMyWaterStatus(c echo.Context) error {
	userID := helper.GetUserID(c)

	result, err := wc.waterService.GetMyWaterStatus(userID)

	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"my water status retrieved successfully",
		result,
	)
}

func (wc *WaterController) GetWaterStatusHistories(c echo.Context) error {
	result, err := wc.waterService.GetWaterStatusHistories()

	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"water histories retrieved successfully",
		result,
	)
}

func (wc *WaterController) SimulateBill(c echo.Context) error {
	var input dto.BillSimulationRequest

	if err := c.Bind(&input); err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	if err := c.Validate(&input); err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	response, err := wc.waterService.SimulateBill(input)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"water bill simulated successfully",
		response,
	)
}

func (wc *WaterController) CreateMeterReading(c echo.Context) error {
	var input dto.CreateMeterReadingInput

	if err := c.Bind(&input); err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	if err := c.Validate(&input); err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	userID := helper.GetUserID(c)

	response, err := wc.waterService.CreateMeterReading(
		userID,
		input,
		fileHeader,
	)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Created(
		c,
		"meter reading submitted successfully",
		response,
	)
}