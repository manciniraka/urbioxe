package repository

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/manciniraka/urbioxe/internal/entity"
	"gorm.io/gorm"
)

type StaffRepository interface {
	CreateStaffTx(tx *gorm.DB, staff *entity.StaffProfile) error
	GetLastEmployeeSequence(departmentID uint, joinDate time.Time) (int, error)
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