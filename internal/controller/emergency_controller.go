package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type EmergencyController struct {
	emergencyService service.EmergencyService
}

func NewEmergencyController(
	emergencyService service.EmergencyService,
) *EmergencyController {
	return &EmergencyController{
		emergencyService: emergencyService,
	}
}

func (ec *EmergencyController) GetAll(c echo.Context) error {

	emergencies, err := ec.emergencyService.GetAll()
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(
		c,
		"emergency contacts fetched successfully",
		emergencies,
	)
}

func (ec *EmergencyController) GetByID(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid emergency contact id",
		)
	}

	emergency, err := ec.emergencyService.GetByID(uint(id))
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(
		c,
		"emergency contact fetched successfully",
		emergency,
	)
}

func (ec *EmergencyController) Create(c echo.Context) error {

	var input service.CreateEmergencyInput

	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(
			c,
			"invalid request body",
		)
	}

	if err := c.Validate(&input); err != nil {
		return helper.BadRequest(
			c,
			err.Error(),
		)
	}

	emergency, err := ec.emergencyService.Create(input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Created(
		c,
		"emergency contact created successfully",
		emergency,
	)
}

func (ec *EmergencyController) Update(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid emergency contact id",
		)
	}

	var input service.UpdateEmergencyInput

	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(
			c,
			"invalid request body",
		)
	}

	if err := c.Validate(&input); err != nil {
		return helper.BadRequest(
			c,
			err.Error(),
		)
	}

	emergency, err := ec.emergencyService.Update(uint(id), input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(
		c,
		"emergency contact updated successfully",
		emergency,
	)
}
