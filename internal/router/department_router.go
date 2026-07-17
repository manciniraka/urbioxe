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

func RegisterDepartmentRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {
	departmentRepository := repository.NewDepartmentRepository(db)
	departmentService := service.NewDepartmentService(departmentRepository)
	departmentController := controller.NewDepartmentController(departmentService)

	department := e.Group("/departments", middleware.AuthMiddleware(cfg))
	department.GET("", departmentController.GetAll, middleware.RequireRoles("citizen", "officer", "department_admin", "super_admin"))
	department.GET("/:id", departmentController.GetByID, middleware.RequireRoles("citizen", "officer", "department_admin", "super_admin"))
	department.POST("", departmentController.Create, middleware.RequireRoles("super_admin"))
	department.PUT("/:id", departmentController.Update, middleware.RequireRoles("super_admin"))
	department.PATCH("/:id/status", departmentController.ToggleStatus, middleware.RequireRoles("super_admin"))

	_ = department
	_ = db
}
