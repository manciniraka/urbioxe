package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterStaffRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	staff := e.Group("/staff")

	_ = staff
	_ = db
}