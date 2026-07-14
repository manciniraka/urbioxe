package router

import (
	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/middleware"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"

	"gorm.io/gorm"
)

func RegisterStaffRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	userRepo := repository.NewUserRepository(db)
	staffRepo := repository.NewStaffRepository(db)

	mailer := mailjet.New(
		mailjet.Config{
			BaseURL:     cfg.MailjetBaseURL,
			APIKey:      cfg.MailjetAPIKey,
			SecretKey:   cfg.MailjetSecretKey,
			SenderEmail: cfg.MailjetSenderEmail,
			SenderName:  cfg.MailjetSenderName,
		},
	)

	staffService := service.NewStaffService(
		db,
		staffRepo,
		userRepo,
		mailer,
	)

	staffController := controller.NewStaffController(staffService)

	staff := e.Group(
		"/staff",
	)

	staff.POST(
		"", 
		staffController.CreateStaff,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(entity.RoleSuperAdmin),
	)

	staff.GET(
		"",
		staffController.GetAllStaff,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(
			entity.RoleSuperAdmin,
			entity.RoleDepartmentAdmin,
		),
	)

	staff.GET(
		"/:id",
		staffController.GetStaffByID,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(
			entity.RoleSuperAdmin,
			entity.RoleDepartmentAdmin,
		),
	)
	staff.PUT(
		"/:id",
		staffController.UpdateStaff,
		middleware.AuthMiddleware(cfg),
		middleware.RequireRoles(
			entity.RoleSuperAdmin,
		),
	)
}