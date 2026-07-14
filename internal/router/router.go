package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/config"
	"gorm.io/gorm"
)

func InitRouter(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {
	// Healthcheck
	e.GET("/", func(c echo.Context) error {
		return c.JSON(
			200,
			echo.Map{
				"message": "urbioxe service is running",
			},
		)
	})

	// Authentication
	RegisterAuthRoutes(e, db, cfg)

	// User
	RegisterUserRoutes(e, db, cfg)

	// Master Data
	RegisterDistrictRoutes(e, db)
	RegisterDepartmentRoutes(e, db)
	RegisterCategoryRoutes(e, db)

	// Staff
	RegisterStaffRoutes(e, db)

	// Reports
	RegisterReportRoutes(e, db, cfg)

	// Public Information
	RegisterRegionalNewsRoutes(e, db)
	RegisterEmergencyContactRoutes(e, db)

	// Weather
	RegisterWeatherRoutes(e, db)
}
