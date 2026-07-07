package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterUserRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	// service
	// controller

	user := e.Group("/users")

	_ = user
	_ = db
}