package service

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type EmergencyService interface {
	GetAll() ([]entity.EmergencyContact, error)
	GetByID(id uint) (*entity.EmergencyContact, error)
	Create(input CreateEmergencyInput) (*entity.EmergencyContact, error)
	Update(id uint, input UpdateEmergencyInput) (*entity.EmergencyContact, error)
}

type emergencyService struct {
	emergencyRepo repository.EmergencyRepository
}

func NewEmergencyService(
	emergencyRepo repository.EmergencyRepository,
) EmergencyService {
	return &emergencyService{
		emergencyRepo: emergencyRepo,
	}
}

type CreateEmergencyInput struct {
	DepartmentID *uint   `json:"department_id"`
	DistrictID   *uint   `json:"district_id"`
	Name         string  `json:"name" validate:"required,min=3,max=100"`
	PhoneNumber  string  `json:"phone_number" validate:"required,min=3,max=15"`
	Description  *string `json:"description"`
	IconURL      *string `json:"icon_url"`
	IsActive     bool    `json:"is_active"`
}

type UpdateEmergencyInput struct {
	DepartmentID *uint   `json:"department_id"`
	DistrictID   *uint   `json:"district_id"`
	Name         string  `json:"name" validate:"required,min=3,max=100"`
	PhoneNumber  string  `json:"phone_number" validate:"required,min=3,max=15"`
	Description  *string `json:"description"`
	IconURL      *string `json:"icon_url"`
	IsActive     bool    `json:"is_active"`
}

func (s *emergencyService) GetAll() ([]entity.EmergencyContact, error) {
	return s.emergencyRepo.GetAll()
}

func (s *emergencyService) GetByID(id uint) (*entity.EmergencyContact, error) {
	emergency, err := s.emergencyRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrEmergencyContactNotFound
		}
		return nil, err
	}

	return emergency, nil
}

func (s *emergencyService) Create(
	input CreateEmergencyInput,
) (*entity.EmergencyContact, error) {

	emergency := entity.EmergencyContact{
		DepartmentID: input.DepartmentID,
		DistrictID:   input.DistrictID,
		Name:         input.Name,
		PhoneNumber:  input.PhoneNumber,
		Description:  input.Description,
		IconURL:      input.IconURL,
		IsActive:     input.IsActive,
	}

	if err := s.emergencyRepo.Create(&emergency); err != nil {
		return nil, err
	}

	return &emergency, nil
}

func (s *emergencyService) Update(
	id uint,
	input UpdateEmergencyInput,
) (*entity.EmergencyContact, error) {

	emergency, err := s.emergencyRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrEmergencyContactNotFound
		}
		return nil, err
	}

	emergency.DepartmentID = input.DepartmentID
	emergency.DistrictID = input.DistrictID
	emergency.Name = input.Name
	emergency.PhoneNumber = input.PhoneNumber
	emergency.Description = input.Description
	emergency.IconURL = input.IconURL
	emergency.IsActive = input.IsActive

	if err := s.emergencyRepo.Update(emergency); err != nil {
		return nil, err
	}

	return emergency, nil
}
