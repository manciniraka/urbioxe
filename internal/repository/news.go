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

func (r *newsRepository) GetAll(ctx context.Context) ([]entity.RegionalNews, error) {
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

func (r *newsRepository) GetByID(ctx context.Context, id int64) (*entity.RegionalNews, error) {
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
		WHERE id = $1
	`
	var news entity.RegionalNews

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&news.ID,
		&news.DepartmentID,
		&news.DistrictID,
		&news.Title,
		&news.Content,
		&news.Category,
		&news.BannerURL,
		&news.TargetScope,
		&news.IsPinned,
		&news.CreatedBy,
		&news.CreatedAt,
		&news.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *newsRepository) Create(ctx context.Context, news *entity.CreateNewsRequest) (*entity.RegionalNews, error) {
	query := `
		INSERT INTO regional_news (
			department_id,
			district_id,
			title,
			content,
			category,
			banner_url,
			target_scope,
			is_pinned
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
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
	`
	var result entity.RegionalNews

	err := r.db.QueryRowContext(
		ctx,
		query,
		news.DepartmentID,
		news.DistrictID,
		news.Title,
		news.Content,
		news.Category,
		news.BannerURL,
		news.TargetScope,
		news.IsPinned,
	).Scan(
		&result.ID,
		&result.DepartmentID,
		&result.DistrictID,
		&result.Title,
		&result.Category,
		&result.BannerURL,
		&result.TargetScope,
		&result.IsPinned,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *newsRepository) Update(ctx context.Context, id int64, news *entity.UpdateNewsRequest) (*entity.RegionalNews, error) {
	query := `
		UPDATE regional_news
		SET
			department_id = $1,
			district_id = $2,
			title = $3,
			content = $4,
			category = $5,
			banner_url = $6,
			target_scope = $7,
			is_pinned = $8,
			updated_at = NOW()
		WHERE id = $9
		RETURNING
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
	`
	var result entity.RegionalNews

	err := r.db.QueryRowContext(
		ctx,
		query,
		news.DepartmentID,
		news.DistrictID,
		news.Title,
		news.Content,
		news.Category,
		news.BannerURL,
		news.TargetScope,
		news.IsPinned,
		id,
	).Scan(
		&result.ID,
		&result.DepartmentID,
		&result.DistrictID,
		&result.Title,
		&result.Content,
		&result.Category,
		&result.BannerURL,
		&result.TargetScope,
		&result.IsPinned,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
