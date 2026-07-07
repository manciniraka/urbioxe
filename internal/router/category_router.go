package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterCategoryRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	categories := e.Group("/categories")

	_ = categories
	_ = db
}