package service

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(userID uint) (*entity.User, error)
	UpdateProfile(userID uint, input UpdateProfileInput) (*entity.User, error)
	ChangePassword(userID uint, input ChangePasswordInput) error
}

type userService struct {
	userRepo repository.UserRepository
	districtRepo repository.DistrictRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	districtRepo repository.DistrictRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		districtRepo: districtRepo,
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

func (s *userService) UpdateProfile(userID uint, input UpdateProfileInput) (*entity.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	if input.HomeDistrictID != nil {
		district, err := s.districtRepo.GetByID(*input.HomeDistrictID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errs.ErrDistrictNotFound
			}
			return nil, err
		}
	
		if !district.IsActive {
			return nil, errs.ErrDistrictInactive
		}
	}

	user.Name = input.Name
	user.PhoneNumber = input.PhoneNumber
	user.HomeDistrictID = input.HomeDistrictID

	err = s.userRepo.UpdateProfile(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *userService) ChangePassword(userID uint, input ChangePasswordInput) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	err = helper.ComparePassword(
		user.Password,
		input.OldPassword,
	)
	if err != nil {
		return errs.ErrInvalidCredential
	}

	if input.OldPassword == input.NewPassword {
		return errs.ErrSamePassword
	}

	hashedPassword, err := helper.HashPassword(
		input.NewPassword,
	)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(
		userID,
		hashedPassword,
	)
	if err != nil {
		return err
	}

	return nil
}
