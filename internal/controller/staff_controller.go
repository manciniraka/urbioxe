package controller

import (
	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type StaffController struct {
	staffService service.StaffService
}

func NewStaffController(
	staffService service.StaffService,
) *StaffController {

	return &StaffController{
		staffService: staffService,
	}
}

func (sc *StaffController) CreateStaff(c echo.Context) error {
	var input service.CreateStaffInput

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

	staff, err := sc.staffService.CreateStaff(
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
		"staff created successfully",
		staff,
	)
}