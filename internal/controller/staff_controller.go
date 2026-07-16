package controller

import (
	"strconv"

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

// CreateStaff godoc
//
//	@Summary		Create staff
//	@Description	Create a new staff account
//	@Tags			Staff
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		service.CreateStaffInput	true	"Staff data"
//	@Success		201		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/staff [post]
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

// GetAllStaff godoc
//
//	@Summary		Get all staff
//	@Description	Retrieve all staff
//	@Tags			Staff
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	helper.Response
//	@Failure		401	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/staff [get]
func (sc *StaffController) GetAllStaff(c echo.Context) error {
	staffs, err := sc.staffService.GetAllStaff()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"staff retrieved successfully",
		staffs,
	)
}

// GetStaffByID godoc
//
//	@Summary		Get staff by ID
//	@Description	Retrieve staff detail
//	@Tags			Staff
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Staff ID"
//	@Success		200	{object}	helper.Response
//	@Failure		400	{object}	helper.ErrorResponse
//	@Failure		404	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/staff/{id} [get]
func (sc *StaffController) GetStaffByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid staff id",
		)
	}

	staff, err := sc.staffService.GetStaffByID(
		uint(id),
	)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"staff retrieved successfully",
		staff,
	)
}

// UpdateStaff godoc
//
//	@Summary		Update staff
//	@Description	Update staff information
//	@Tags			Staff
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Staff ID"
//	@Param			request	body		service.UpdateStaffInput	true	"Staff data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/staff/{id} [put]
func (sc *StaffController) UpdateStaff(c echo.Context) error {
	id, err := strconv.Atoi(
		c.Param("id"),
	)
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid staff id",
		)
	}

	var input service.UpdateStaffInput

	if err := c.Bind(
		&input,
	); err != nil {
		return helper.BadRequest(
			c,
			"invalid request body",
		)
	}

	if err := c.Validate(
		&input,
	); err != nil {
		return helper.BadRequest(
			c,
			err.Error(),
		)
	}

	staff, err := sc.staffService.UpdateStaff(
		uint(id),
		input,
	)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"staff updated successfully",
		staff,
	)
}
