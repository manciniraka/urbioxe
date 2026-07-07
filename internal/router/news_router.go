package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRegionalNewsRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	news := e.Group("/news")

	_ = news
	_ = db
}