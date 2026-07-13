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
