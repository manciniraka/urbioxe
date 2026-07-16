package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	FindAll(isOnlyActive bool) ([]entity.Department, error)
	FindByID(id uint) (*entity.Department, error)
	Create(department *entity.Department) error
	Update(department *entity.Department) error
	UpdateStatus(id uint, isActive bool) error
	CheckCodeOrNameExists(code string, name string, excludeID uint) (bool, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func (dr *departmentRepository) FindAll(isOnlyActive bool) ([]entity.Department, error) {
	var departments []entity.Department

	query := dr.db.Model(&entity.Department{})

	if isOnlyActive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("name ASC").Find(&departments).Error
	if err != nil {
		return nil, err
	}

	return departments, nil
}

func (dr *departmentRepository) FindByID(id uint) (*entity.Department, error) {
	var department entity.Department

	err := dr.db.Model(&entity.Department{}).
		First(&department, id).Error

	if err != nil {
		return nil, err
	}

	return &department, nil
}

func (dr *departmentRepository) Create(department *entity.Department) error {
	return dr.db.Create(department).Error
}

func (dr *departmentRepository) Update(department *entity.Department) error {
	return dr.db.Save(department).Error
}

func (dr *departmentRepository) UpdateStatus(id uint, isActive bool) error {
	return dr.db.Model(&entity.Department{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

func (dr *departmentRepository) CheckCodeOrNameExists(code string, name string, excludeID uint) (bool, error) {
	var count int64

	query := dr.db.Model(&entity.Department{}).
		Where("LOWER(code) = LOWER(?) OR LOWER(name) = LOWER(?)", code, name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
