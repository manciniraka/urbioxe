package service

import (
	"errors"

	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(input RegisterInput) (*entity.User, error)
	Login(input LoginInput) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
	mailer   *mailjet.Client
}

func NewAuthService(
	userRepo repository.UserRepository,
	cfg *config.Config,
	mailer *mailjet.Client,
) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
		mailer:   mailer,
	}
}

type RegisterInput struct {
	NIK            string `json:"nik" validate:"required,len=16,numeric"`
	HomeDistrictID *uint  `json:"home_district_id"`
	Name           string `json:"name" validate:"required,min=3,max=100"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	PhoneNumber    string `json:"phone_number" validate:"required,min=9,max=15,numeric"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (s *authService) Register(input RegisterInput) (*entity.User, error) {
	existingUser, err := s.userRepo.FindByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, errs.ErrEmailRegistered
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existingUser, err = s.userRepo.FindByNIK(input.NIK)
	if err == nil && existingUser != nil {
		return nil, errs.ErrNIKRegistered
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		HomeDistrictID: input.HomeDistrictID,
		NIK:            input.NIK,
		Name:           input.Name,
		Email:          input.Email,
		Password:       hashedPassword,
		PhoneNumber:    input.PhoneNumber,
		Role:           entity.RoleCitizen,
	}

	err = s.userRepo.Register(&user)
	if err := s.mailer.SendWelcomeEmail(
		user.Name,
		user.Email,
	); err != nil {
		return nil, err
	}

	user.Password = ""

	return &user, nil
}

// Validates user credentials to login
func (s *authService) Login(input LoginInput) (string, error) {
	user, err := s.userRepo.FindByEmail(
		input.Email,
	)
	if err != nil {
		return "", errs.ErrInvalidCredential
	}

	err = helper.ComparePassword(
		user.Password,
		input.Password,
	)
	if err != nil {
		return "", errs.ErrInvalidCredential
	}

	token, err := helper.GenerateJWT(
		user.ID,
		string(user.Role),
		s.cfg.JWTSecret,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
