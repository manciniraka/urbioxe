package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterDistrictRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	district := e.Group("/districts")

	_ = district
	_ = db
}