package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(isOnlyActive bool) ([]entity.Category, error)
	FindByID(id int64) (*entity.Category, error)
	Create(category *entity.Category) error
	Update(category *entity.Category) error
	UpdateStatus(id int64, isActive bool) error
	CheckNameExistsInDepartment(departmentID int64, name string, excludeID int64) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (cr *categoryRepository) FindAll(isOnlyActive bool) ([]entity.Category, error) {
	var categories []entity.Category

	query := cr.db.Model(&entity.Category{}).Preload("Department")

	if isOnlyActive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (cr *categoryRepository) FindByID(id int64) (*entity.Category, error) {
	var category entity.Category

	err := cr.db.Model(&entity.Category{}).
		Preload("Department").
		First(&category, id).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (cr *categoryRepository) Create(category *entity.Category) error {
	return cr.db.Create(category).Error
}

func (cr *categoryRepository) Update(category *entity.Category) error {
	return cr.db.Save(category).Error
}

func (cr *categoryRepository) UpdateStatus(id int64, isActive bool) error {
	return cr.db.Model(&entity.Category{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

func (cr *categoryRepository) CheckNameExistsInDepartment(departmentID int64, name string, excludeID int64) (bool, error) {
	var count int64

	query := cr.db.Model(&entity.Category{}).
		Where("department_id = ? AND LOWER(name) = LOWER(?)", departmentID, name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
