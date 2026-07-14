package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
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

	department := e.Group("/departments")
	department.GET("", departmentController.GetAll)
	department.GET("/:id", departmentController.GetByID)
	department.POST("", departmentController.Create)
	department.PUT("/:id", departmentController.Update)
	department.PATCH("/:id", departmentController.ToggleStatus)

	_ = department
	_ = db
}
