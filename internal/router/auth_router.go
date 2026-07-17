package router

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/controller"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/service"
)

func RegisterAuthRoutes(
	e *echo.Echo,
	db *gorm.DB,
	cfg *config.Config,
) {

	userRepo := repository.NewUserRepository(db)
	mailer := mailjet.New(
		mailjet.Config{
			BaseURL:     cfg.MailjetBaseURL,
			APIKey:      cfg.MailjetAPIKey,
			SecretKey:   cfg.MailjetSecretKey,
			SenderEmail: cfg.MailjetSenderEmail,
			SenderName:  cfg.MailjetSenderName,
		},
	)
	authService := service.NewAuthService(userRepo, cfg, mailer)
	authController := controller.NewAuthController(authService)

	e.POST("/register", authController.Register)

	e.POST("/login", authController.Login)

	_ = cfg
}
