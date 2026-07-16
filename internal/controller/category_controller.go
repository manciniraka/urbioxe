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

// GetAll godoc
//
//	@Summary		Get all categories
//	@Description	Get all categories
//	@Tags			Category
//	@Produce		json
//	@Param			active_only	query		bool	false	"Show only active categories"
//	@Success		200			{object}	helper.Response
//	@Failure		500			{object}	helper.ErrorResponse
//	@Router			/categories [get]
func (cc *CategoryController) GetAll(c echo.Context) error {
	activeOnlyParam := c.QueryParam("active_only")
	isOnlyActive := activeOnlyParam == "true"

	categories, err := cc.svc.GetAllCategories(isOnlyActive)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get all categories", categories)
}

// GetByID godoc
//
//	@Summary		Get category by ID
//	@Description	Retrieve category detail
//	@Tags			Category
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	helper.Response
//	@Failure		400	{object}	helper.ErrorResponse
//	@Failure		404	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/categories/{id} [get]
func (cc *CategoryController) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "category id is not valid")
	}

	category, err := cc.svc.GetCategoryByID(uint(id))
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "succes get category", category)
}

// Create godoc
//
//	@Summary		Create category
//	@Description	Create a new category
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		service.CreateCategoryInput	true	"Category data"
//	@Success		201		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/categories [post]
func (cc *CategoryController) Create(c echo.Context) error {
	role := helper.GetUserRole(c)

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

// Update godoc
//
//	@Summary		Update category
//	@Description	Update category
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Category ID"
//	@Param			request	body		service.UpdateCategoryInput	true	"Category data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/categories/{id} [put]
func (cc *CategoryController) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "category id is not valid")
	}

	role := helper.GetUserRole(c)

	var input service.UpdateCategoryInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body is not valid")
	}

	category, err := cc.svc.UpdateCategory(uint(id), role, input)
	if err != nil {
		if err.Error() == "department_id required" ||
			err.Error() == "category name required" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success update category", category)
}

// ToggleStatus godoc
//
//	@Summary		Toggle category status
//	@Description	Enable or disable category
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"Category ID"
//	@Param			request	body		service.ToggleCategoryStatusInput	true	"Status data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/categories/{id}/status [patch]
func (cc *CategoryController) ToggleStatus(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "category id is not valid")
	}

	role := helper.GetUserRole(c)

	var input service.ToggleCategoryStatusInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body is not valid")
	}

	err = cc.svc.ToggleCategoryStatus(uint(id), role, input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success update status category", nil)
}
