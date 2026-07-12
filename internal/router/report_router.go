package router

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/external/cloudinary"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"
	"gorm.io/gorm"
)

func RegisterReportRoutes(
	e *echo.Echo,
	db *gorm.DB,
) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cldService := cloudinary.NewCloudinaryService(
		cloudName,
		apiKey,
		apiSecret,
	)
	reportRepo := repository.NewReportRepository(db)
	reportSvc := service.NewReportService(reportRepo, cldService)
	reportCtrl := controller.NewReportController(reportSvc)

	report := e.Group("/reports")
	report.POST("", reportCtrl.Create)
	report.GET("", reportCtrl.GetAll)
	report.GET("/:id", reportCtrl.GetByID)

	_ = report
	_ = db
}
