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

func RegisterRegionalNewsRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {
	newsRepo := repository.NewNewsRepository(db)
	newsService := service.NewNewsService(newsRepo)
	newsController := controller.NewNewsController(newsService)

	news := e.Group("/news")

	news.GET("", newsController.GetAll)
	news.GET("/:id", newsController.GetByID)
	news.POST("", newsController.Create,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles("department_admin", "super_admin"))
	news.PUT("/:id", newsController.Update,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles("department_admin", "super_admin"))
	news.DELETE("/:id", newsController.Delete,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles("department_admin", "super_admin"))
}
