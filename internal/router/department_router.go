package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterDepartmentRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	department := e.Group("/departments")

	_ = department
	_ = db
}