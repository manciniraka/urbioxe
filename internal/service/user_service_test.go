package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	mockUserRepo     *mocks.MockUserRepository
	mockDistrictRepo *mocks.MockDistrictRepository
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.ctrl)
	suite.mockDistrictRepo = mocks.NewMockDistrictRepository(suite.ctrl)
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *UserServiceTestSuite) TestGetProfile() {
	userID := uint(1)

	suite.Run("Success - Get profile details", func() {
		suite.SetupTest()
		existingUser := &entity.User{
			ID:       userID,
			Name:     "Rian",
			Email:    "rian@mail.com",
			Password: "encrypted_secret_password",
		}

		suite.mockUserRepo.EXPECT().GetByID(userID).Return(existingUser, nil)

		svc := service.NewUserService(suite.mockUserRepo, suite.mockDistrictRepo)
		res, err := svc.GetProfile(userID)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Empty(suite.T(), res.Password)
	})

	suite.Run("Fail - User Not Found", func() {
		suite.SetupTest()
		suite.mockUserRepo.EXPECT().GetByID(userID).Return(nil, gorm.ErrRecordNotFound)

		svc := service.NewUserService(suite.mockUserRepo, suite.mockDistrictRepo)
		res, err := svc.GetProfile(userID)

		assert.ErrorIs(suite.T(), err, errs.ErrUserNotFound)
		assert.Nil(suite.T(), res)
	})
}

func (suite *UserServiceTestSuite) TestUpdateProfile() {
	userID := uint(1)
	activeDistrictID := uint(3)
	inactiveDistrictID := uint(4)
	notFoundDistrictID := uint(99)

	tests := []struct {
		name     string
		input    service.UpdateProfileInput
		mockFn   func()
		verifyFn func(res *entity.User, err error)
	}{
		{
			name: "Success - Update profile with active district",
			input: service.UpdateProfileInput{
				Name:           "Rian Diperbarui",
				PhoneNumber:    "089999999",
				HomeDistrictID: &activeDistrictID,
			},
			mockFn: func() {
				suite.mockUserRepo.EXPECT().GetByID(userID).Return(&entity.User{ID: userID, Name: "Rian Lama"}, nil)
				suite.mockDistrictRepo.EXPECT().GetByID(activeDistrictID).Return(&entity.District{ID: activeDistrictID, IsActive: true}, nil)
				suite.mockUserRepo.EXPECT().UpdateProfile(gomock.Any()).Return(nil)
			},
			verifyFn: func(res *entity.User, err error) {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), res)
				assert.Equal(suite.T(), "Rian Diperbarui", res.Name)
				assert.Empty(suite.T(), res.Password)
			},
		},
		{
			name: "Success - Update profile with nil district",
			input: service.UpdateProfileInput{
				Name:           "Rian Tanpa District",
				PhoneNumber:    "089999999",
				HomeDistrictID: nil,
			},
			mockFn: func() {
				suite.mockUserRepo.EXPECT().GetByID(userID).Return(&entity.User{ID: userID, Name: "Rian Lama"}, nil)
				suite.mockUserRepo.EXPECT().UpdateProfile(gomock.Any()).Return(nil)
			},
			verifyFn: func(res *entity.User, err error) {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), res)
				assert.Equal(suite.T(), "Rian Tanpa District", res.Name)
			},
		},
		{
			name: "Fail - User Not Found",
			input: service.UpdateProfileInput{
				Name:        "Rian Ghosting",
				PhoneNumber: "089999999",
			},
			mockFn: func() {
				suite.mockUserRepo.EXPECT().GetByID(userID).Return(nil, gorm.ErrRecordNotFound)
			},
			verifyFn: func(res *entity.User, err error) {
				assert.ErrorIs(suite.T(), err, errs.ErrUserNotFound)
				assert.Nil(suite.T(), res)
			},
		},
		{
			name: "Fail - District Not Found",
			input: service.UpdateProfileInput{
				Name:           "Rian Wrong District",
				PhoneNumber:    "089999999",
				HomeDistrictID: &notFoundDistrictID,
			},
			mockFn: func() {
				suite.mockUserRepo.EXPECT().GetByID(userID).Return(&entity.User{ID: userID}, nil)
				suite.mockDistrictRepo.EXPECT().GetByID(notFoundDistrictID).Return(nil, gorm.ErrRecordNotFound)
			},
			verifyFn: func(res *entity.User, err error) {
				assert.ErrorIs(suite.T(), err, errs.ErrDistrictNotFound)
				assert.Nil(suite.T(), res)
			},
		},
		{
			name: "Fail - District Inactive",
			input: service.UpdateProfileInput{
				Name:           "Rian Inactive District",
				PhoneNumber:    "089999999",
				HomeDistrictID: &inactiveDistrictID,
			},
			mockFn: func() {
				suite.mockUserRepo.EXPECT().GetByID(userID).Return(&entity.User{ID: userID}, nil)
				suite.mockDistrictRepo.EXPECT().GetByID(inactiveDistrictID).Return(&entity.District{ID: inactiveDistrictID, IsActive: false}, nil)
			},
			verifyFn: func(res *entity.User, err error) {
				assert.ErrorIs(suite.T(), err, errs.ErrDistrictInactive)
				assert.Nil(suite.T(), res)
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockFn()

			svc := service.NewUserService(suite.mockUserRepo, suite.mockDistrictRepo)
			res, err := svc.UpdateProfile(userID, tc.input)

			tc.verifyFn(res, err)
		})
	}
}

func (suite *UserServiceTestSuite) TestChangePassword() {
	userID := uint(1)
	input := service.ChangePasswordInput{
		OldPassword: "password123",
		NewPassword: "password456",
	}

	hashedOldPassword, err := helper.HashPassword(input.OldPassword)
	if err != nil {
		suite.T().Fatal("failed to generate dummy hash for testing:", err)
	}

	suite.Run("Success - Password changed successfully", func() {
		suite.SetupTest()
		existingUser := &entity.User{
			ID:       userID,
			Password: hashedOldPassword,
		}

		suite.mockUserRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		suite.mockUserRepo.EXPECT().UpdatePassword(userID, gomock.Any()).Return(nil)

		svc := service.NewUserService(suite.mockUserRepo, suite.mockDistrictRepo)
		err := svc.ChangePassword(userID, input)

		assert.Nil(suite.T(), err)
	})

	suite.Run("Fail - Same Password Detected", func() {
		suite.SetupTest()
		sameInput := service.ChangePasswordInput{
			OldPassword: "password123",
			NewPassword: "password123",
		}
		existingUser := &entity.User{
			ID:       userID,
			Password: hashedOldPassword,
		}

		suite.mockUserRepo.EXPECT().GetByID(userID).Return(existingUser, nil)

		svc := service.NewUserService(suite.mockUserRepo, suite.mockDistrictRepo)
		err := svc.ChangePassword(userID, sameInput)

		assert.ErrorIs(suite.T(), err, errs.ErrSamePassword)
	})
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
