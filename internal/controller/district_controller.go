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