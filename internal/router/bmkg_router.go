package router

import (
	externalbmkg "github.com/manciniraka/urbioxe/external/bmkg"

	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/middleware"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/scheduler"
	"github.com/manciniraka/urbioxe/internal/service"

	"gorm.io/gorm"
)

func RegisterBMKGRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	bmkgClient := externalbmkg.New(
		externalbmkg.Config{
			BaseURL:          cfg.BMKGBaseURL,
			ForecastEndpoint: cfg.BMKGForecastEndpoint,
		},
	)

	bmkgRepo := repository.NewBMKGRepository(db)
	districtRepo := repository.NewDistrictRepository(db)
	userRepo := repository.NewUserRepository(db)
	bmkgService := service.NewBMKGService(
		db,
		bmkgClient,
		bmkgRepo,
		districtRepo,
		userRepo,
	)
	bmkgController := controller.NewBMKGController(bmkgService)
	bmkgScheduler := scheduler.NewBMKGScheduler(bmkgService)

	bmkgScheduler.Start()

	bmkg := e.Group("/bmkg")

	bmkg.GET("", bmkgController.GetAllWeather)
	bmkg.GET(
		"/me",
		bmkgController.GetMyWeather,
		middleware.AuthMiddleware(cfg),
	)
	bmkg.GET("/:district_id",
		bmkgController.GetWeatherDetailsByDistrictID,
		middleware.AuthMiddleware(cfg),
	)
	bmkg.POST(
		"/sync",
		bmkgController.SyncForecasts,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(entity.RoleSuperAdmin),
	)
}
