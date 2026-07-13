package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/external/openweather"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/middleware"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/scheduler"
	"github.com/manciniraka/urbioxe/internal/service"
	"gorm.io/gorm"
)

func RegisterWeatherRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	weatherRepo := repository.NewWeatherRepository(db)
	userRepo := repository.NewUserRepository(db)

	openWeather := openweather.New(
		openweather.Config{
			BaseURL: cfg.OpenWeatherBaseURL,
			APIKey:  cfg.OpenWeatherAPIKey,
		},
	)

	weatherService := service.NewWeatherService(
		weatherRepo,
		userRepo,
		openWeather,
	)

	weatherController := controller.NewWeatherController(weatherService)

	weatherScheduler := scheduler.NewWeatherScheduler(
		weatherService,
	)

	weatherScheduler.Start()

	weather := e.Group("/weather")

	weather.GET("", weatherController.GetAllWeather)
	weather.GET(
		"/me",
		weatherController.GetMyWeather,
		middleware.AuthMiddleware(cfg),
	)
	weather.GET("/:district_id", weatherController.GetWeatherByDistrictID)
	weather.POST(
		"/sync",
		weatherController.SyncWeather,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(entity.RoleSuperAdmin),
	)
}
