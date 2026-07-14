package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/repository"
	"github.com/manciniraka/urbioxe/internal/validator"
)

type NewsService interface {
	GetAll(ctx context.Context) ([]entity.RegionalNews, error)
	GetByID(ctx context.Context, id int64) (*entity.RegionalNews, error)
	Create(ctx context.Context, news *entity.CreateNewsRequest) (*entity.RegionalNews, error)
	Update(ctx context.Context, id int64, news *entity.UpdateNewsRequest) (*entity.RegionalNews, error)
	Delete(ctx context.Context, id int64) error
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

func (s *newsService) GetAll(
	ctx context.Context,
) ([]entity.RegionalNews, error) {
	return s.newsRepository.GetAll(ctx)
}

func (s *newsService) GetByID(
	ctx context.Context,
	id int64,
) (*entity.RegionalNews, error) {
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
	ctx context.Context,
	news *entity.CreateNewsRequest,
) (*entity.RegionalNews, error) {
	if err := validator.ValidateCreateNews(news); err != nil {
		return nil, err
	}

	return s.newsRepository.Create(ctx, news)
}

func (s *newsService) Update(
	ctx context.Context,
	id int64,
	news *entity.UpdateNewsRequest,
) (*entity.RegionalNews, error) {
	return s.newsRepository.Update(ctx, id, news)
}

func (s *newsService) Delete(
	ctx context.Context,
	id int64,
) error {
	err := s.newsRepository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}
	return nil
}
