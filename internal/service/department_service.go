package service

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/repository"
)

type CreateDepartmentInput struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateDepartmentInput struct {
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
	CreateDepartment(role string, input CreateDepartmentInput) (*entity.Department, error)
	UpdateDepartment(id uint, role string, input UpdateDepartmentInput) (*entity.Department, error)
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
	return &entity.Department{}, nil
}

func (ds *departmentService) CreateDepartment(role string, input CreateDepartmentInput) (*entity.Department, error) {
	return &entity.Department{}, nil
}

func (ds *departmentService) UpdateDepartment(id uint, role string, input UpdateDepartmentInput) (*entity.Department, error) {
	return &entity.Department{}, nil
}

func (ds *departmentService) ToggleDepartmentStatus(id uint, role string, input ToggleDepartmentStatusInput) error {
	return nil
}
