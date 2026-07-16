package service

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/repository"
)

type DistrictService interface {
	GetAllDistrict() ([]entity.District, error)
	GetDistrictByID(id uint) (*entity.District, error)
	CreateDistrict(input CreateDistrictInput) (*entity.District, error)
	UpdateDistrict(id uint, input UpdateDistrictInput) (*entity.District, error)
}

type districtService struct {
	districtRepo repository.DistrictRepository
}

func NewDistrictService(
	districtRepo repository.DistrictRepository,
) DistrictService {

	return &districtService{
		districtRepo: districtRepo,
	}
}

type CreateDistrictInput struct {
	Name string `json:"name" validate:"required"`
}

type UpdateDistrictInput struct {
	IsActive bool `json:"is_active"`
}

func (ds *districtService) GetAllDistrict() ([]entity.District, error) {
	return ds.districtRepo.GetAll()
}

func (ds *districtService) GetDistrictByID(id uint) (*entity.District, error) {
	return ds.districtRepo.GetByID(id)
}

func (ds *districtService) CreateDistrict(input CreateDistrictInput) (*entity.District, error) {
	existing, err := ds.districtRepo.GetByName(
		input.Name,
	)

	if err == nil && existing != nil {
		return nil, errs.ErrDistrictAlreadyExists
	}

	if err != nil &&
		!errors.Is(err, errs.ErrDistrictNotFound) {
		return nil, err
	}

	district := entity.District{
		Name: input.Name,
	}

	if err := ds.districtRepo.Create(
		&district,
	); err != nil {
		return nil, err
	}

	return &district, nil
}

func (ds *districtService) UpdateDistrict(id uint, input UpdateDistrictInput) (*entity.District, error) {
	district, err := ds.districtRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	district.IsActive = input.IsActive

	if err := ds.districtRepo.Update(district); err != nil {
		return nil, err
	}

	return district, nil
}
