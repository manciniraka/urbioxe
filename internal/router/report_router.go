package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterReportRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	report := e.Group("/reports")

	_ = report
	_ = db
}