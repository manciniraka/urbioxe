package mocks

import (
	"context"

	"github.com/manciniraka/urbioxe/internal/entity"
)

type NewsRepositoryMock struct {
	GetAllFunc func(ctx context.Context) ([]entity.RegionalNews, error)

	GetByIDFunc func(ctx context.Context, id int64) (*entity.RegionalNews, error)

	CreateFunc func(ctx context.Context, news *entity.CreateNewsRequest) (*entity.RegionalNews, error)

	UpdateFunc func(ctx context.Context, id int64, news *entity.UpdateNewsRequest) (*entity.RegionalNews, error)

	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *NewsRepositoryMock) GetAll(ctx context.Context) ([]entity.RegionalNews, error) {
	return m.GetAllFunc(ctx)
}

func (m *NewsRepositoryMock) GetByID(ctx context.Context, id int64) (*entity.RegionalNews, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *NewsRepositoryMock) Create(ctx context.Context, news *entity.CreateNewsRequest) (*entity.RegionalNews, error) {
	return m.CreateFunc(ctx, news)
}

func (m *NewsRepositoryMock) Update(ctx context.Context, id int64, news *entity.UpdateNewsRequest) (*entity.RegionalNews, error) {
	return m.UpdateFunc(ctx, id, news)
}

func (m *NewsRepositoryMock) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
