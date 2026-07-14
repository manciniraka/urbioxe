package repository

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"gorm.io/gorm"
)

type StaffRepository interface {
	CreateStaffTx(tx *gorm.DB, staff *entity.StaffProfile) error
	GetLastEmployeeSequence(departmentID uint, joinDate time.Time) (int, error)

	FindStaffByUserID(userID uint) (*entity.StaffProfile, error)
	GetAll() ([]entity.StaffProfile, error)
	GetByID(id uint) (*entity.StaffProfile, error)
	UpdateStaffTx(tx *gorm.DB, staff *entity.StaffProfile) error
}

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(
	db *gorm.DB,
) StaffRepository {
	return &staffRepository{
		db: db,
	}
}

func (sr *staffRepository) CreateStaffTx(tx *gorm.DB, staff *entity.StaffProfile) error {
	return tx.Create(staff).Error
}

func (sr *staffRepository) GetLastEmployeeSequence(departmentID uint, joinDate time.Time) (int, error) {
	var staff entity.StaffProfile

	err := sr.db.
		Where(
			"department_id = ?",
			departmentID,
		).
		Where(
			"EXTRACT(YEAR FROM join_date) = ?",
			joinDate.Year(),
		).
		Order("employee_number DESC").
		First(&staff).Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return 1, nil
		}

		return 0, err
	}

	parts := strings.Split(
		staff.EmployeeNumber,
		"-",
	)

	if len(parts) != 3 {
		return 1, nil
	}

	lastSequence, err := strconv.Atoi(
		parts[2],
	)
	if err != nil {
		return 0, err
	}

	return lastSequence + 1, nil
}

func (sr *staffRepository) FindStaffByUserID(userID uint) (*entity.StaffProfile, error) {
	var staff entity.StaffProfile

	err := sr.db.
		Preload("User").
		Preload("Department").
		Where(
			"user_id = ?",
			userID,
		).
		First(&staff).Error

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errs.ErrStaffNotFound
		}

		return nil, err
	}

	return &staff, nil
}

func (sr *staffRepository) GetAll() ([]entity.StaffProfile, error) {
	var staffs []entity.StaffProfile

	err := sr.db.
		Preload("User").
		Preload("Department").
		Order("employee_number ASC").
		Find(&staffs).Error

	if err != nil {
		return nil, err
	}

	return staffs, nil
}

func (sr *staffRepository) GetByID(id uint) (*entity.StaffProfile, error) {
	var staff entity.StaffProfile

	err := sr.db.
		Preload("User").
		Preload("Department").
		First(
			&staff,
			id,
		).Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errs.ErrStaffNotFound
		}

		return nil, err
	}

	return &staff, nil
}

func (sr *staffRepository) UpdateStaffTx(tx *gorm.DB, staff *entity.StaffProfile) error{
	return tx.
		Model(&entity.StaffProfile{}).
		Where("id = ?", staff.ID).
		Updates(map[string]any{
			"department_id": staff.DepartmentID,
			"position":      staff.Position,
			"is_active":     staff.IsActive,
		}).Error
}