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

// GetAll godoc
//
//	@Summary		Get all emergency contacts
//	@Description	Retrieve all emergency contacts
//	@Tags			Emergency
//	@Produce		json
//	@Success		200	{object}	helper.Response
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/emergency-contacts [get]
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

// GetByID godoc
//
//	@Summary		Get emergency contact by ID
//	@Description	Retrieve emergency contact detail
//	@Tags			Emergency
//	@Produce		json
//	@Param			id	path		int	true	"Emergency Contact ID"
//	@Success		200	{object}	helper.Response
//	@Failure		400	{object}	helper.ErrorResponse
//	@Failure		404	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/emergency-contacts/{id} [get]
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

// Create godoc
//
//	@Summary		Create emergency contact
//	@Description	Create a new emergency contact
//	@Tags			Emergency
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//  @Param 			request body 		service.CreateEmergencyInput true "Emergency contact data"
//	@Success		201		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/emergency-contacts [post]
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

// Update godoc
//
//	@Summary		Update emergency contact
//	@Description	Update emergency contact
//	@Tags			Emergency
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"Emergency Contact ID"
// 	@Param 			request body 		service.UpdateEmergencyInput true "Emergency contact data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/emergency-contacts/{id} [put]
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
