package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/config"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	// TODO:
	// Repository
	// Service
	// Controller

	auth := e.Group("")

	_ = auth
	_ = db
	_ = cfg
}