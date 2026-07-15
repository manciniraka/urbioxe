package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/middleware"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"
	"gorm.io/gorm"
)

func RegisterEmergencyContactRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
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
	emergency.POST("", emergencyController.Create,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles("department_admin", "super_admin"))
	emergency.PUT("/:id", emergencyController.Update,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles("department_admin", "super_admin"))

	_ = emergency
	_ = db
}
