package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type DistrictController struct {
	districtService service.DistrictService
}

func NewDistrictController(
	districtService service.DistrictService,
) *DistrictController {

	return &DistrictController{
		districtService: districtService,
	}
}

// GetAllDistrict godoc
//
//	@Summary		Get all districts
//	@Description	Retrieve all districts
//	@Tags			District
//	@Produce		json
//	@Success		200	{object}	helper.Response
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/districts [get]
func (dc *DistrictController) GetAllDistrict(c echo.Context) error {
	districts, err := dc.districtService.GetAllDistrict()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"district retrieved successfully",
		districts,
	)
}

// GetDistrictByID godoc
//
//	@Summary		Get district by ID
//	@Description	Retrieve district detail
//	@Tags			District
//	@Produce		json
//	@Param			id	path		int	true	"District ID"
//	@Success		200	{object}	helper.Response
//	@Failure		400	{object}	helper.ErrorResponse
//	@Failure		404	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/districts/{id} [get]
func (dc *DistrictController) GetDistrictByID(c echo.Context) error {
	id, err := strconv.Atoi(
		c.Param("id"),
	)
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid district id",
		)
	}

	district, err := dc.districtService.GetDistrictByID(
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
		"district retrieved successfully",
		district,
	)
}

// CreateDistrict godoc
//
//	@Summary		Create district
//	@Description	Create a new district
//	@Tags			District
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		service.CreateDistrictInput	true	"District data"
//	@Success		201		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/districts [post]
func (dc *DistrictController) CreateDistrict(c echo.Context) error {
	var input service.CreateDistrictInput

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

	district, err := dc.districtService.CreateDistrict(
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
		"district created successfully",
		district,
	)
}

// UpdateDistrict godoc
//
//	@Summary		Update district
//	@Description	Update district
//	@Tags			District
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"District ID"
//	@Param			request	body		service.UpdateDistrictInput	true	"District data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/districts/{id} [put]
func (dc *DistrictController) UpdateDistrict(c echo.Context) error {
	id, err := strconv.Atoi(
		c.Param("id"),
	)
	if err != nil {

		return helper.BadRequest(
			c,
			"invalid district id",
		)
	}

	var input service.UpdateDistrictInput

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

	district, err := dc.districtService.UpdateDistrict(
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
		"district updated successfully",
		district,
	)
}
