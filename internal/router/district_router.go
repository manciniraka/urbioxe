package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/middleware"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"
	"gorm.io/gorm"
)

func RegisterDistrictRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	districtRepo := repository.NewDistrictRepository(db)
	districtService := service.NewDistrictService(districtRepo)
	districtController := controller.NewDistrictController(districtService)

	district := e.Group(
		"/districts",
	)

	district.GET("", districtController.GetAllDistrict)
	district.GET("/:id", districtController.GetDistrictByID)
	district.POST(
		"",
		districtController.CreateDistrict,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(
			entity.RoleSuperAdmin,
		),
	)
	district.PUT(
		"/:id",
		districtController.UpdateDistrict,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(
			entity.RoleSuperAdmin,
		),
	)
}
