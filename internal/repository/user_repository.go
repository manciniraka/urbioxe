package repository

import (
	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(user *entity.User) error
	RegisterUserTx(tx *gorm.DB, user *entity.User) error

	FindByEmail(email string) (*entity.User, error)
	FindByNIK(nik string) (*entity.User, error)

	GetByID(id uint) (*entity.User, error)
	UpdateProfile(user *entity.User) error
	UpdatePassword(id uint, password string) error
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

func (ur *userRepository) RegisterUserTx(tx *gorm.DB, user *entity.User) error {
	return tx.Create(user).Error
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

func (ur *userRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User

	err := ur.db.
		First(&user, id).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) UpdateProfile(user *entity.User) error {
	return ur.db.
		Save(user).
		Error
}

func (ur *userRepository) UpdatePassword(id uint, password string) error {
	return ur.db.
		Model(&entity.User{}).
		Where("id = ?", id).
		Update("password", password).
		Error
}
