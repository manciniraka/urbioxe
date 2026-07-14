package service

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/validator"
)

type CreateNewsInput struct {
	DepartmentID *int64              `json:"department_id"`
	DistrictID   *int64              `json:"district_id,omitempty"`
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	Category     entity.NewsCategory `json:"category"`
	BannerURL    *string             `json:"banner_url"`
	TargetScope  entity.NewsScope    `json:"target_scope"`
	IsPinned     bool                `json:"is_pinned"`
}

type UpdateNewsInput struct {
	DepartmentID *int64              `json:"department_id"`
	DistrictID   *int64              `json:"district_id,omitempty"`
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	Category     entity.NewsCategory `json:"category"`
	BannerURL    *string             `json:"banner_url"`
	TargetScope  entity.NewsScope    `json:"target_scope"`
	IsPinned     bool                `json:"is_pinned"`
}

type NewsService interface {
	GetAll() ([]entity.RegionalNews, error)
	GetByID(id int64) (*entity.RegionalNews, error)
	Create(
		createdBy int64,
		input CreateNewsInput,
	) (*entity.RegionalNews, error)
	Update(
		id int64,
		input UpdateNewsInput,
	) (*entity.RegionalNews, error)
	Delete(id int64) error
}

type newsService struct {
	newsRepository repository.NewsRepository
}

func NewNewsService(
	newsRepository repository.NewsRepository,
) NewsService {
	return &newsService{
		newsRepository: newsRepository,
	}
}

func (s *newsService) GetAll() ([]entity.RegionalNews, error) {
	return s.newsRepository.GetAll()
}

func (s *newsService) GetByID(
	id int64,
) (*entity.RegionalNews, error) {
	return s.newsRepository.GetByID(id)
}

func (s *newsService) Create(
	createdBy int64,
	input CreateNewsInput,
) (*entity.RegionalNews, error) {
	news := &entity.RegionalNews{
		DepartmentID: input.DepartmentID,
		DistrictID:   input.DistrictID,
		Title:        input.Title,
		Content:      input.Content,
		Category:     input.Category,
		BannerURL:    input.BannerURL,
		TargetScope:  input.TargetScope,
		IsPinned:     input.IsPinned,
		CreatedBy:    &createdBy,
	}

	if err := validator.ValidateCreateNews(&entity.CreateNewsRequest{
		DepartmentID: input.DepartmentID,
		DistrictID:   input.DistrictID,
		Title:        input.Title,
		Content:      input.Content,
		Category:     input.Category,
		BannerURL:    input.BannerURL,
		TargetScope:  input.TargetScope,
		IsPinned:     input.IsPinned,
	}); err != nil {
		return nil, err
	}

	if err := s.newsRepository.Create(news); err != nil {
		return nil, err
	}

	return news, nil
}

func (s *newsService) Update(
	id int64,
	input UpdateNewsInput,
) (*entity.RegionalNews, error) {
	news, err := s.newsRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	news.DepartmentID = input.DepartmentID
	news.DistrictID = input.DistrictID
	news.Title = input.Title
	news.Content = input.Content
	news.Category = input.Category
	news.BannerURL = input.BannerURL
	news.TargetScope = input.TargetScope
	news.IsPinned = input.IsPinned

	if err := s.newsRepository.Update(news); err != nil {
		return nil, err
	}

	return news, nil
}

func (s *newsService) Delete(
	id int64,
) error {
	return s.newsRepository.Delete(id)
}
