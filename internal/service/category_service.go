package service

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type CategoryService interface {
	GetAllCategories(isOnlyActive bool) ([]entity.Category, error)
	GetCategoryByID(id int64) (*entity.Category, error)
	CreateCategory(role string, input CreateCategoryInput) (*entity.Category, error)
	UpdateCategory(id int64, role string, input UpdateCategoryInput) (*entity.Category, error)
	ToggleCategoryStatus(id int64, role string, input ToggleCategoryStatusInput) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

type CreateCategoryInput struct {
	DepartmentID int64  `json:"department_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type UpdateCategoryInput struct {
	DepartmentID int64  `json:"department_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type ToggleCategoryStatusInput struct {
	IsActive bool `json:"is_active"`
}

func (cs *categoryService) GetAllCategories(isOnlyActive bool) ([]entity.Category, error) {
	return cs.repo.FindAll(isOnlyActive)
}

func (cs *categoryService) GetCategoryByID(id int64) (*entity.Category, error) {
	category, err := cs.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (cs *categoryService) CreateCategory(role string, input CreateCategoryInput) (*entity.Category, error) {
	if role != "super_admin" && role != "department_admin" {
		return nil, errs.ErrForbidden
	}

	if input.DepartmentID == 0 {
		return nil, errors.New("department_id required")
	}
	if input.Name == "" {
		return nil, errors.New("category name required")
	}

	exists, err := cs.repo.CheckNameExistsInDepartment(input.DepartmentID, input.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.ErrCategoryAlreadyExists
	}

	category := entity.Category{
		DepartmentID: input.DepartmentID,
		Name:         input.Name,
		Description:  input.Description,
		IsActive:     true,
	}

	err = cs.repo.Create(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (cs *categoryService) UpdateCategory(id int64, role string, input UpdateCategoryInput) (*entity.Category, error) {
	if role != "super_admin" && role != "department_admin" {
		return nil, errs.ErrForbidden
	}

	category, err := cs.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, err
	}

	if input.DepartmentID == 0 {
		return nil, errors.New("department_id required")
	}
	if input.Name == "" {
		return nil, errors.New("category name required")
	}

	exists, err := cs.repo.CheckNameExistsInDepartment(input.DepartmentID, input.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.ErrCategoryAlreadyExists
	}

	category.DepartmentID = input.DepartmentID
	category.Name = input.Name
	category.Description = input.Description

	err = cs.repo.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *categoryService) ToggleCategoryStatus(id int64, role string, input ToggleCategoryStatusInput) error {
	if role != "super_admin" && role != "department_admin" {
		return errs.ErrForbidden
	}

	_, err := cs.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrCategoryNotFound
		}
		return err
	}

	return cs.repo.UpdateStatus(id, input.IsActive)
}
