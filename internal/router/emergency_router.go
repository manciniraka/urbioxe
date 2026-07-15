package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"
	"gorm.io/gorm"
)

func RegisterEmergencyContactRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	// repository
	emergencyRepo := repository.NewEmergencyRepository(db)
	// service
	emergencyService := service.NewEmergencyService(emergencyRepo)
	// controller
	emergencyController := controller.NewEmergencyController(emergencyService)

	emergency := e.Group("/emergency-contacts")
	// Public
	emergency.GET("", emergencyController.GetAll)
	emergency.GET("/:id", emergencyController.GetByID)

	// Department Admin / Super Admin
	emergency.POST("", emergencyController.Create)
	emergency.PUT("/:id", emergencyController.Update)

	_ = emergency
	_ = db
}
