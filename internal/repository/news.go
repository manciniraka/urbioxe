package repository

import (
	"context"
	"database/sql"

	"github.com/manciniraka/urbioxe/internal/entity"
)

type NewsRepository interface {
	GetAll(ctx context.Context) ([]entity.RegionalNews, error)
	GetByID(ctx context.Context, id int64) (*entity.RegionalNews, error)
	Create(ctx context.Context, news *entity.CreateNewsRequest) (*entity.RegionalNews, error)
	Update(ctx context.Context, id int64, news *entity.UpdateNewsRequest) (*entity.RegionalNews, error)
	Delete(ctx context.Context, id int64) error
}

type newsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) *newsRepository {
	return &newsRepository{
		db: db,
	}
}

func (r *newsRepository) GetAll(ctx context.Context, ) ([]entity.RegionalNews, error) {
	query := `
		SELECT
			id,
			department_id,
			district_id,
			title,
			content,
			category,
			banner_url,
			target_scope,
			is_pinned,
			created_by,
			created_at,
			updated_at
		FROM regional_news
		ORDER BY is_pinned DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	news := make([]entity.RegionalNews, 0)

	for rows.Next() {
		var item entity.RegionalNews

		err := rows.Scan(
			&item.ID,
			&item.DepartmentID,
			&item.DistrictID,
			&item.Title,
			&item.Content,
			&item.Category,
			&item.BannerURL,
			&item.TargetScope,
			&item.IsPinned,
			&item.CreatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		news = append(news, item)
	}

	return news, nil
}