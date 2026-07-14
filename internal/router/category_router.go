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

func RegisterCategoryRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	// TODO:
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryController := controller.NewCategoryController(categoryService)

	categories := e.Group("/categories", middleware.AuthMiddleware(cfg))
	categories.GET("", categoryController.GetAll, middleware.RequireRoles("citizen", "officer", "department_admin", "super_admin"))
	categories.GET("/:id", categoryController.GetByID, middleware.RequireRoles("citizen", "officer", "department_admin", "super_admin"))
	categories.POST("", categoryController.Create, middleware.RequireRoles("super_admin"))
	categories.PUT("/:id", categoryController.Update, middleware.RequireRoles("super_admin"))
	categories.PATCH("/:id/status", categoryController.ToggleStatus, middleware.RequireRoles("super_admin"))

	_ = categories
	_ = db
}
