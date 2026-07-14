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

func (dc *DepartmentController) GetAll(c echo.Context) error {
	activeOnlyParam := c.QueryParam("active_only")
	isOnlyActive := activeOnlyParam == "true"

	departments, err := dc.svc.GetAllDepartments(isOnlyActive)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get all departments", departments)
}

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

func (dc *DepartmentController) Create(c echo.Context) error {
	role := helper.GetUserRole(c)
	// test
	if role == "" {
		role = "super_admin"
	}

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

func (dc *DepartmentController) Update(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "department id is not valid")
	}

	role := helper.GetUserRole(c)
	// not login
	if role == "" {
		role = "super_admin"
	}

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

func (dc *DepartmentController) ToggleStatus(c echo.Context) error {
	return nil
}
