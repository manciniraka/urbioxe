package controller

import (
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
