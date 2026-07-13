package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type CategoryController struct {
	svc service.CategoryService
}

func NewCategoryController(svc service.CategoryService) *CategoryController {
	return &CategoryController{svc: svc}
}

func (cc *CategoryController) GetAll(c echo.Context) error {
	activeOnlyParam := c.QueryParam("active_only")
	isOnlyActive := activeOnlyParam == "true"

	categories, err := cc.svc.GetAllCategories(isOnlyActive)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get all categories", categories)
}

func (cc *CategoryController) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "category id is not valid")
	}

	category, err := cc.svc.GetCategoryByID(id)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "succes get category", category)
}

func (cc *CategoryController) Create(c echo.Context) error {
	role, _ := c.Get("role").(string)
	// test without login
	if role == "" {
		role = "department_admin"
	}

	var input service.CreateCategoryInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	category, err := cc.svc.CreateCategory(role, input)
	if err != nil {
		if err.Error() == "department_id required" || err.Error() == "category name required" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Created(c, "success create category", category)
}

func (cc *CategoryController) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "category id is not valid")
	}

	role, _ := c.Get("role").(string)
	if role == "" {
		role = "department_admin"
	}

	var input service.UpdateCategoryInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body is not valid")
	}

	category, err := cc.svc.UpdateCategory(id, role, input)
	if err != nil {
		if err.Error() == "department_id required" ||
			err.Error() == "category name required" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success update category", category)
}
