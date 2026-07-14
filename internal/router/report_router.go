package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/external/cloudinary"
	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"
	"gorm.io/gorm"
)

func RegisterReportRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {
	cloudName := cfg.CloudinaryCloudName
	apiKey := cfg.CloudinaryAPIKey
	apiSecret := cfg.CloudinaryAPISecret

	cldService := cloudinary.NewCloudinaryService(
		cloudName,
		apiKey,
		apiSecret,
	)

	mailer := mailjet.New(
		mailjet.Config{
			BaseURL:     cfg.MailjetBaseURL,
			APIKey:      cfg.MailjetAPIKey,
			SecretKey:   cfg.MailjetSecretKey,
			SenderEmail: cfg.MailjetSenderEmail,
			SenderName:  cfg.MailjetSenderName,
		},
	)

	reportRepo := repository.NewReportRepository(db)
	reportSvc := service.NewReportService(reportRepo, cldService, mailer)
	reportCtrl := controller.NewReportController(reportSvc)

	report := e.Group("/reports")
	report.POST("", reportCtrl.Create)
	report.GET("", reportCtrl.GetAll)
	report.GET("/:id", reportCtrl.GetByID)
	report.PUT("/:id", reportCtrl.Update)
	report.PATCH("/:id/assign", reportCtrl.Assign)
	report.PATCH("/:id/start", reportCtrl.Start)
	report.PATCH("/:id/resolve", reportCtrl.Resolve)
	report.PATCH("/:id/reject", reportCtrl.Reject)

	_ = report
	_ = db
}
