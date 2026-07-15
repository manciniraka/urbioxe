package repository

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"

	"gorm.io/gorm"
)

type DistrictRepository interface {
	GetAll() ([]entity.District, error)

	GetByID(id uint) (*entity.District, error)
	GetByName(name string) (*entity.District, error)

	Create(district *entity.District) error
	Update(district *entity.District) error
}

type districtRepository struct {
	db *gorm.DB
}

func NewDistrictRepository(
	db *gorm.DB,
) DistrictRepository {

	return &districtRepository{
		db: db,
	}
}

func (dr *districtRepository) GetAll() ([]entity.District, error) {
	var districts []entity.District

	err := dr.db.
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&districts).Error

	if err != nil {
		return nil, err
	}

	return districts, nil
}

func (dr *districtRepository) GetByID(id uint) (*entity.District, error){
	var district entity.District

	err := dr.db.
		First(
			&district,
			id,
		).Error

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errs.ErrDistrictNotFound
		}

		return nil, err
	}

	return &district, nil
}

func (dr *districtRepository) GetByName(name string) (*entity.District, error){
	var district entity.District

	err := dr.db.
		Where(
			"name = ?",
			name,
		).
		First(&district).Error

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errs.ErrDistrictNotFound
		}

		return nil, err
	}

	return &district, nil
}

func (dr *districtRepository) Create(district *entity.District) error{
	return dr.db.Create(district).Error
}

func (dr *districtRepository) Update(district *entity.District) error{
	return dr.db.
		Model(&entity.District{}).
		Where("id = ?",	district.ID).
		Updates(
			map[string]any{
				"is_active": district.IsActive,
			},
		).Error
}