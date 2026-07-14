package router

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/external/cloudinary"
	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/middleware"
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

	staffRepo := repository.NewStaffRepository(db)

	reportRepo := repository.NewReportRepository(db)
	reportSvc := service.NewReportService(reportRepo, staffRepo, cldService, mailer)
	reportCtrl := controller.NewReportController(reportSvc)

	report := e.Group("/reports", middleware.AuthMiddleware(cfg))
	report.POST("", reportCtrl.Create, middleware.RequireRoles("citizen"))
	report.PUT("/:id", reportCtrl.Update, middleware.RequireRoles("citizen"))
	report.GET("", reportCtrl.GetAll, middleware.RequireRoles("citizen", "officer", "department_admin", "super_admin"))
	report.GET("/:id", reportCtrl.GetByID, middleware.RequireRoles("citizen", "officer", "department_admin", "super_admin"))
	report.PATCH("/:id/assign", reportCtrl.Assign, middleware.RequireRoles("department_admin", "super_admin"))
	report.PATCH("/:id/start", reportCtrl.Start, middleware.RequireRoles("officer"))
	report.PATCH("/:id/resolve", reportCtrl.Resolve, middleware.RequireRoles("officer"))
	report.PATCH("/:id/reject", reportCtrl.Reject, middleware.RequireRoles("department_admin", "super_admin"))
	report.PATCH("/:id/verify", reportCtrl.Verify, middleware.RequireRoles("department_admin", "super_admin"))
	report.PATCH("/:id/priority", reportCtrl.UpdatePriority, middleware.RequireRoles("department_admin", "super_admin"))

	_ = report
	_ = db
}
