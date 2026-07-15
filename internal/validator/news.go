package validator

import (
	"errors"

	"github.com/manciniraka/urbioxe/internal/entity"
)

func ValidateCreateNews(news *entity.CreateNewsRequest) error {
	if news.Title == "" {
		return errors.New("title is required")
	}

	if news.Content == "" {
		return errors.New("Content is required")
	}

	switch news.Category {
	case entity.NewsCategoryAnnouncement,
		entity.NewsCategoryNews,
		entity.NewsCategoryEvent,
		entity.NewsCategoryEmergency:
	default:
		return errors.New("invalid news category")
	}

	switch news.TargetScope {
	case entity.NewsScopeGlobal:
	case entity.NewsScopeDistrict:
		if news.DistrictID == nil {
			return errors.New("district_id is required for district scope")
		}
	default:
		return errors.New("invalid news target scope")
	}

	return nil
}
