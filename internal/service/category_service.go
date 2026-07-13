package service

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/repository"
)

type CategoryService interface {
	GetAllCategories(isOnlyActive bool) ([]entity.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (cs *categoryService) GetAllCategories(isOnlyActive bool) ([]entity.Category, error) {
	return cs.repo.FindAll(isOnlyActive)
}
