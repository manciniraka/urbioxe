package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"
	"gorm.io/gorm"
)

func RegisterCategoryRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {

	// TODO:
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryController := controller.NewCategoryController(categoryService)

	categories := e.Group("/categories")
	categories.GET("", categoryController.GetAll)

	_ = categories
	_ = db
}
