package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockUserRepo *mocks.MockUserRepository
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.ctrl)
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

		svc := service.NewUserService(suite.mockUserRepo)
		res, err := svc.GetProfile(userID)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Empty(suite.T(), res.Password)
	})

	suite.Run("Fail - User Not Found", func() {
		suite.SetupTest()
		suite.mockUserRepo.EXPECT().GetByID(userID).Return(nil, errors.New("not found"))

		svc := service.NewUserService(suite.mockUserRepo)
		res, err := svc.GetProfile(userID)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *UserServiceTestSuite) TestUpdateProfile() {
	userID := uint(1)
	districtID := uint(3)
	input := service.UpdateProfileInput{
		Name:           "Rian Diperbarui",
		PhoneNumber:    "089999999",
		HomeDistrictID: &districtID,
	}

	suite.Run("Success - Update identity profile", func() {
		suite.SetupTest()
		existingUser := &entity.User{
			ID:       userID,
			Name:     "Rian Lama",
			Password: "secret_password",
		}

		suite.mockUserRepo.EXPECT().GetByID(userID).Return(existingUser, nil)
		suite.mockUserRepo.EXPECT().UpdateProfile(gomock.Any()).Return(nil)

		svc := service.NewUserService(suite.mockUserRepo)
		res, err := svc.UpdateProfile(userID, input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), "Rian Diperbarui", res.Name)
		assert.Empty(suite.T(), res.Password)
	})
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

		svc := service.NewUserService(suite.mockUserRepo)
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

		svc := service.NewUserService(suite.mockUserRepo)
		err := svc.ChangePassword(userID, sameInput)

		assert.NotNil(suite.T(), err)
	})
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
