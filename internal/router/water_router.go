package router

import (
	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/external/cloudinary"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/middleware"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"

	"gorm.io/gorm"
)

func RegisterWaterRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {
	waterRepo := repository.NewWaterRepository(db)
	districtRepo := repository.NewDistrictRepository(db)
	userRepo := repository.NewUserRepository(db)

	cloudinaryService := cloudinary.NewCloudinaryService(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)

	meterReadingRepo := repository.NewMeterReadingRepository(db)

	waterService := service.NewWaterService(
		waterRepo,
		districtRepo,
		userRepo,
		meterReadingRepo,
		cloudinaryService,
	)
	waterController := controller.NewWaterController(waterService)

	water := e.Group("/water")

	water.GET("", waterController.GetAllWaterStatus)
	water.GET(
		"/me",
		waterController.GetMyWaterStatus,
		middleware.AuthMiddleware(cfg),
	)
	water.GET("/:district_id", waterController.GetWaterStatusByDistrictID)
	water.POST(
		"",
		waterController.CreateWaterStatus,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(
			entity.RoleDepartmentAdmin,
			entity.RoleSuperAdmin,
		),
	)
	water.GET(
		"/histories",
		waterController.GetWaterStatusHistories,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(
			entity.RoleOfficer,
			entity.RoleDepartmentAdmin,
			entity.RoleSuperAdmin,
		),
	)

	water.POST("/bill/simulate", waterController.SimulateBill)

	water.POST(
		"/meter-reading",
		waterController.CreateMeterReading,
		middleware.AuthMiddleware(cfg),
	)
}