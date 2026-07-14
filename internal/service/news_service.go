package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/validator"
)

type CreateNewsInput struct {
	DepartmentID *int64  				`json:"department_id"`
	DistrictID   *int64 				`json:"district_id,omitempty"`
	Title        string 				`json:"title"`
	Content      string 				`json:"content"`
	Category     entity.NewsCategory 	`json:"category"`
	BannerURL    *string 				`json:"banner_url"`
	TargetScope  entity.NewsScope 		`json:"target_scope"`
	IsPinned     bool   				`json:"is_pinned"`
}

type UpdateNewsInput struct {
	DepartmentID *int64  				`json:"department_id"`
	DistrictID   *int64 				`json:"district_id,omitempty"`
	Title        string 				`json:"title"`
	Content      string 				`json:"content"`
	Category     entity.NewsCategory 	`json:"category"`
	BannerURL    *string 				`json:"banner_url"`
	TargetScope  entity.NewsScope 		`json:"target_scope"`
	IsPinned     bool   				`json:"is_pinned"`
}

type NewsService interface {
	GetAll() ([]entity.RegionalNews, error)
	GetByID(id int64) (*entity.RegionalNews, error)
	Create(createdBy int64, input CreateNewsInput) (*entity.RegionalNews, error)
	Update(id int64, news *entity.UpdateNewsRequest) (*entity.RegionalNews, error)
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
	ctx := context.Background()

	return s.newsRepository.GetAll(ctx)
}

func (s *newsService) GetByID(
	id int64,
) (*entity.RegionalNews, error) {
	ctx := context.Background()
	news, err := s.newsRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return news, nil
}

func (s *newsService) Create(
	createdBy int64,
	input CreateNewsInput,
) (*entity.RegionalNews, error) {
	news := &entity.CreateNewsRequest{
		DepartmentID: input.DepartmentID,
		DistrictID:   input.DistrictID,
		Title:        input.Title,
		Content:      input.Content,
		Category:     input.Category,
		BannerURL:    input.BannerURL,
		TargetScope:  input.TargetScope,
		IsPinned:     input.IsPinned,
	}

	if err := validator.ValidateCreateNews(news); err != nil {
		return nil, err
	}

	ctx := context.Background()

	return s.newsRepository.Create(
		ctx,
		createdBy,
		news,
	)
}

func (s *newsService) Update(
	id int64,
	news *entity.UpdateNewsRequest,
) (*entity.RegionalNews, error) {
	ctx := context.Background()
	return s.newsRepository.Update(ctx, id, news)
}

func (s *newsService) Delete(
	id int64,
) error {
	ctx := context.Background()
	
	err := s.newsRepository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}
	return nil
}
