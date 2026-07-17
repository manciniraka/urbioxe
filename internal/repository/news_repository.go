package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type NewsRepository interface {
	GetAll() ([]entity.RegionalNews, error)
	GetByID(id int64) (*entity.RegionalNews, error)
	Create(news *entity.RegionalNews) error
	Update(news *entity.RegionalNews) error
	Delete(id int64) error
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(
	db *gorm.DB,
) NewsRepository {
	return &newsRepository{
		db: db,
	}
}

func (nr *newsRepository) GetAll() ([]entity.RegionalNews, error) {
	var news []entity.RegionalNews

	err := nr.db.
		Order("is_pinned DESC").
		Order("created_at DESC").
		Find(&news).
		Error

	if err != nil {
		return nil, err
	}

	return news, nil
}

func (nr *newsRepository) GetByID(
	id int64,
) (*entity.RegionalNews, error) {
	var news entity.RegionalNews

	err := nr.db.
		First(&news, id).
		Error

	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (nr *newsRepository) Create(
	news *entity.RegionalNews,
) error {
	return nr.db.
		Create(news).
		Error
}

func (nr *newsRepository) Update(
	news *entity.RegionalNews,
) error {
	return nr.db.
		Save(news).
		Error
}

func (nr *newsRepository) Delete(
	id int64,
) error {
	result := nr.db.
		Delete(
			&entity.RegionalNews{},
			id,
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
