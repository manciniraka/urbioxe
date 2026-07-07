package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByNIK(nik string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Register(user *entity.User) error {
	return ur.db.Create(user).Error
}

func (ur *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := ur.db.
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) FindByNIK(nik string) (*entity.User, error) {
	var user entity.User

	err := ur.db.
		Where("nik = ?", nik).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}