package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type EmergencyRepository interface {
	GetAll() ([]entity.EmergencyContact, error)
	GetByID(id uint) (*entity.EmergencyContact, error)
	Create(emergency *entity.EmergencyContact) error
	Update(emergency *entity.EmergencyContact) error
}

type emergencyRepository struct {
	db *gorm.DB
}

func NewEmergencyRepository(db *gorm.DB) EmergencyRepository {
	return &emergencyRepository{
		db: db,
	}
}

func (er *emergencyRepository) GetAll() ([]entity.EmergencyContact, error) {
	var emergencies []entity.EmergencyContact

	err := er.db.Find(&emergencies).Error
	if err != nil {
		return nil, err
	}

	return emergencies, nil
}

func (er *emergencyRepository) GetByID(id uint) (*entity.EmergencyContact, error) {
	var emergency entity.EmergencyContact

	err := er.db.
		Where("id = ?", id).
		First(&emergency).
		Error

	if err != nil {
		return nil, err
	}

	return &emergency, nil
}

func (er *emergencyRepository) Create(emergency *entity.EmergencyContact) error {
	return er.db.Create(emergency).Error
}

func (er *emergencyRepository) Update(emergency *entity.EmergencyContact) error {
	return er.db.Save(emergency).Error
}
