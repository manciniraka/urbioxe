package service

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type DepartmentInput struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ToggleDepartmentStatusInput struct {
	IsActive bool `json:"is_active"`
}

type DepartmentService interface {
	GetAllDepartments(isOnlyActive bool) ([]entity.Department, error)
	GetDepartmentByID(id uint) (*entity.Department, error)
	CreateDepartment(role string, input DepartmentInput) (*entity.Department, error)
	UpdateDepartment(id uint, role string, input DepartmentInput) (*entity.Department, error)
	ToggleDepartmentStatus(id uint, role string, input ToggleDepartmentStatusInput) error
}

type departmentService struct {
	repo repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{
		repo: repo,
	}
}

func (ds *departmentService) GetAllDepartments(isOnlyActive bool) ([]entity.Department, error) {
	return ds.repo.FindAll(isOnlyActive)
}

func (ds *departmentService) GetDepartmentByID(id uint) (*entity.Department, error) {
	dept, err := ds.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrDepartmentNotFound
		}
		return nil, err
	}
	return dept, nil
}

func (ds *departmentService) CreateDepartment(role string, input DepartmentInput) (*entity.Department, error) {
	if role != "super_admin" {
		return nil, errs.ErrForbidden
	}

	if input.Code == "" {
		return nil, errors.New("departement code required")
	}
	if input.Name == "" {
		return nil, errors.New("departement name required")
	}

	exists, err := ds.repo.CheckCodeOrNameExists(input.Code, input.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("department code or name already used")
	}

	dept := entity.Department{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		IsActive:    true,
	}

	err = ds.repo.Create(&dept)
	if err != nil {
		return nil, err
	}

	return &dept, nil
}

func (ds *departmentService) UpdateDepartment(id uint, role string, input DepartmentInput) (*entity.Department, error) {
	if role != "super_admin" {
		return nil, errs.ErrForbidden
	}

	dept, err := ds.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrDepartmentNotFound
		}
		return nil, err
	}

	if input.Code == "" {
		return nil, errors.New("departement code required")
	}
	if input.Name == "" {
		return nil, errors.New("departement name required")
	}

	exists, err := ds.repo.CheckCodeOrNameExists(input.Code, input.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("department code or name already used")
	}

	dept.Code = input.Code
	dept.Name = input.Name
	dept.Description = input.Description

	err = ds.repo.Update(dept)
	if err != nil {
		return nil, err
	}

	return dept, nil
}

func (ds *departmentService) ToggleDepartmentStatus(id uint, role string, input ToggleDepartmentStatusInput) error {
	return nil
}
