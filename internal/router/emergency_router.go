package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterEmergencyContactRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	emergency := e.Group("/emergency-contacts")

	_ = emergency
	_ = db
}
