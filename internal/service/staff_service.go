package service

import (
	"errors"
	"time"

	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/constant"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/logger"
	"github.com/manciniraka/urbioxe/internal/repository"
	"gorm.io/gorm"
)

type StaffService interface {
	CreateStaff(input CreateStaffInput) (*entity.StaffProfile, error)
}

type staffService struct {
	db        *gorm.DB
	staffRepo repository.StaffRepository
	userRepo  repository.UserRepository
	mailer    *mailjet.Client
}

func NewStaffService(
	db *gorm.DB,
	staffRepo repository.StaffRepository,
	userRepo repository.UserRepository,
	mailer *mailjet.Client,
) StaffService {
	return &staffService{
		db:        db,
		staffRepo: staffRepo,
		userRepo:  userRepo,
		mailer:    mailer,
	}
}

type CreateStaffInput struct {
	NIK          string          `json:"nik" validate:"required,len=16"`
	Name         string          `json:"name" validate:"required"`
	Email        string          `json:"email" validate:"required,email"`
	PhoneNumber  string          `json:"phone_number"`
	DepartmentID uint            `json:"department_id" validate:"required"`
	Position     entity.StaffPosition `json:"position" validate:"required"`
	JoinDate string `json:"join_date" validate:"required"`
}

func (ss *staffService) CreateStaff(input CreateStaffInput) (*entity.StaffProfile, error) {
	var userRole entity.UserRole

	switch input.Position {
		case entity.PositionFieldOfficer:
			userRole = entity.RoleOfficer
		case entity.PositionDepartmentAdmin:
			userRole = entity.RoleDepartmentAdmin
		case entity.PositionSupervisor:
			userRole = entity.RoleSuperAdmin

		default:
			return nil, errs.ErrInvalidStaffPosition
	}

	existingUser, err := ss.userRepo.FindByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, errs.ErrEmailRegistered
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	existingUser, err = ss.userRepo.FindByNIK(input.NIK)
	if err == nil && existingUser != nil {
		return nil, errs.ErrNIKRegistered
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var department entity.Department

	err = ss.db.
		First( // TODO(INTEGRATION): Replace direct query with DepartmentRepository after Department module has been merged.
			&department,
			input.DepartmentID,
		).Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, errs.ErrDepartmentNotFound
		}

		return nil, err
	}

	joinDate, err := time.Parse(
	    "2006-01-02",
	    input.JoinDate,
	)
	if err != nil {
	    return nil, errs.ErrBadRequest
	}

	lastSequence, err := ss.staffRepo.GetLastEmployeeSequence(
	    input.DepartmentID,
	    joinDate,
	)
	if err != nil {
		return nil, err
	}

	employeeNumber := helper.GenerateEmployeeNumber(
	    department.Code,
	    joinDate,
	    lastSequence,
	)

	temporaryPassword, err := helper.GenerateTemporaryPassword(department.Code)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := helper.HashPassword(temporaryPassword)
	if err != nil {
		return nil, err
	}

	var (
		user  entity.User
		staff entity.StaffProfile
	)

	err = ss.db.Transaction(
		func(tx *gorm.DB) error {
			user = entity.User{
				NIK:         input.NIK,
				Name:        input.Name,
				Email:       input.Email,
				Password:    hashedPassword,
				PhoneNumber: input.PhoneNumber,
				Role:        userRole,
			}

			if err := ss.userRepo.RegisterUserTx(
				tx,
				&user,
			); err != nil {
				return err
			}

			staff = entity.StaffProfile{
				UserID:         user.ID,
				DepartmentID:   input.DepartmentID,
				EmployeeNumber: employeeNumber,
				Position:       input.Position,
				JoinDate:       input.JoinDate,
				IsActive:       true,
			}

			if err := ss.staffRepo.CreateStaffTx(
				tx,
				&staff,
			); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	if err := ss.mailer.SendWelcomeStaffEmail(
		user.Name,
		user.Email,
		department.Name,
		staff.EmployeeNumber,
		temporaryPassword,
	); err != nil {
		logger.Log.Error(
			"failed to send welcome staff email",
			"tag", constant.LogTagMailjet,
			"email", user.Email,
			"error", err,
		)
	}

	user.Password = ""

	staff.User = &user
	staff.Department = &department

	return &staff, nil
}