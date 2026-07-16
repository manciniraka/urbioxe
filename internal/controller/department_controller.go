package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type DepartmentController struct {
	svc service.DepartmentService
}

func NewDepartmentController(svc service.DepartmentService) *DepartmentController {
	return &DepartmentController{svc: svc}
}

// GetAll godoc
//
//	@Summary		Get all departments
//	@Description	Retrieve all departments
//	@Tags			Department
//	@Produce		json
//	@Param			active_only	query		bool	false	"Show active departments only"
//	@Success		200			{object}	helper.Response
//	@Failure		500			{object}	helper.ErrorResponse
//	@Router			/departments [get]
func (dc *DepartmentController) GetAll(c echo.Context) error {
	activeOnlyParam := c.QueryParam("active_only")
	isOnlyActive := activeOnlyParam == "true"

	departments, err := dc.svc.GetAllDepartments(isOnlyActive)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get all departments", departments)
}

// GetByID godoc
//
//	@Summary		Get department by ID
//	@Description	Retrieve department detail
//	@Tags			Department
//	@Produce		json
//	@Param			id	path		int	true	"Department ID"
//	@Success		200	{object}	helper.Response
//	@Failure		400	{object}	helper.ErrorResponse
//	@Failure		404	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/departments/{id} [get]
func (dc *DepartmentController) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "department id is not valid")
	}

	dept, err := dc.svc.GetDepartmentByID(uint(id64))
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get department detaail", dept)
}

// Create godoc
//
//	@Summary		Create department
//	@Description	Create a new department
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		service.DepartmentInput	true	"Department data"
//	@Success		201		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/departments [post]
func (dc *DepartmentController) Create(c echo.Context) error {
	role := helper.GetUserRole(c)

	var input service.DepartmentInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body is not valid")
	}

	dept, err := dc.svc.CreateDepartment(role, input)
	if err != nil {
		if err.Error() == "departement code required" ||
			err.Error() == "departement name required" ||
			err.Error() == "department code or name already used" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Created(c, "department created", dept)
}

// Update godoc
//
//	@Summary		Update department
//	@Description	Update department
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Department ID"
//	@Param			request	body		service.DepartmentInput	true	"Department data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/departments/{id} [put]
func (dc *DepartmentController) Update(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "department id is not valid")
	}

	role := helper.GetUserRole(c)

	var input service.DepartmentInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body is not valid")
	}

	dept, err := dc.svc.UpdateDepartment(uint(id64), role, input)
	if err != nil {
		if err.Error() == "departement code required" ||
			err.Error() == "departement name required" ||
			err.Error() == "department code or name already used" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success update department", dept)
}

// ToggleStatus godoc
//
//	@Summary		Toggle department status
//	@Description	Enable or disable department
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"Department ID"
//	@Param			request	body		service.ToggleDepartmentStatusInput	true	"Status data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/departments/{id}/status [patch]
func (dc *DepartmentController) ToggleStatus(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "department id is not valid")
	}

	role := helper.GetUserRole(c)

	var input service.ToggleDepartmentStatusInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body is not valid")
	}

	err = dc.svc.ToggleDepartmentStatus(uint(id64), role, input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success update department status", nil)
}
