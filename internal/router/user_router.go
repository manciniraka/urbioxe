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

func RegisterUserRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	userRepo := repository.NewUserRepository(db)
	districtRepo := repository.NewDistrictRepository(db)
	userService := service.NewUserService(userRepo, districtRepo)
	userController := controller.NewUserController(userService)

	users := e.Group(
		"/users",
		middleware.AuthMiddleware(cfg),
	)

	users.GET("/profile", userController.GetProfile)
	users.PUT("/profile", userController.UpdateProfile)
	users.PUT("/change-password", userController.ChangePassword)
}
