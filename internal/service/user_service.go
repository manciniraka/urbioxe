package service

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/repository"
)

type UserService interface {
	GetProfile(userID uint) (*entity.User, error)
	UpdateProfile(userID uint, input UpdateProfileInput) (*entity.User, error)
	ChangePassword(userID uint, input ChangePasswordInput) (*entity.User, error)
}

type userService struct{
	userRepo repository.UserRepository
}

func NewUserService(
	userRepo repository.UserRepository,
) UserService{
	return &userService{
		userRepo: userRepo,
	}
}

type UpdateProfileInput struct {
	Name           string `json:"name" validate:"required,min=3,max=100"`
	PhoneNumber    string `json:"phone_number" validate:"required,min=9,max=15,numeric"`
	HomeDistrictID *uint  `json:"home_district_id"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (s *userService) GetProfile(userID uint) (*entity.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	user.Password = ""

	return user, nil
}