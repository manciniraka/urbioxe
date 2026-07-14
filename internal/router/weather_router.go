package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterWeatherRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	weather := e.Group("/weathers")

	_ = weather
	_ = db
}