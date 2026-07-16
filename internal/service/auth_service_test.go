package service_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func uintPtr(n uint) *uint {
	return &n
}

type AuthServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockUserRepo *mocks.MockUserRepository
	mailer       *mailjet.Client
	cfg          *config.Config
	mockServer   *httptest.Server
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.ctrl)

	suite.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))

	suite.mailer = mailjet.New(
		mailjet.Config{
			BaseURL:     suite.mockServer.URL,
			APIKey:      "dummy_key",
			SecretKey:   "dummy_secret",
			SenderEmail: "noreply@urbioxe.com",
			SenderName:  "Urbioxe System",
		},
	)

	suite.cfg = &config.Config{
		JWTSecret: "supersecretjwtkey",
	}
}

func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.mockServer.Close()
	suite.ctrl.Finish()
}

func (suite *AuthServiceTestSuite) TestRegister() {
	input := service.RegisterInput{
		HomeDistrictID: uintPtr(1),
		NIK:            "3201010101010001",
		Name:           "Rian Citizen",
		Email:          "rian@mail.com",
		Password:       "password123",
		PhoneNumber:    "08123456789",
	}

	tests := []struct {
		name          string
		mockFn        func()
		expectedError string
	}{
		{
			name: "Success - Register citizen",
			mockFn: func() {
				suite.mockUserRepo.EXPECT().
					FindByEmail(input.Email).
					Return(nil, gorm.ErrRecordNotFound)

				suite.mockUserRepo.EXPECT().
					FindByNIK(input.NIK).
					Return(nil, gorm.ErrRecordNotFound)

				suite.mockUserRepo.EXPECT().
					Register(gomock.Any()).
					DoAndReturn(func(user *entity.User) error {
						user.ID = 1
						return nil
					})

			},
			expectedError: "",
		},
		{
			name: "Fail - Email already registered",
			mockFn: func() {
				existingUser := &entity.User{ID: 1, Email: input.Email}
				suite.mockUserRepo.EXPECT().
					FindByEmail(input.Email).
					Return(existingUser, nil)
			},
			expectedError: "",
		},
		{
			name: "Fail - NIK already registered",
			mockFn: func() {
				suite.mockUserRepo.EXPECT().
					FindByEmail(input.Email).
					Return(nil, gorm.ErrRecordNotFound)

				existingUser := &entity.User{ID: 1, NIK: input.NIK}
				suite.mockUserRepo.EXPECT().
					FindByNIK(input.NIK).
					Return(existingUser, nil)
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockFn()

			svc := service.NewAuthService(suite.mockUserRepo, suite.cfg, suite.mailer)
			result, err := svc.Register(input)

			if tc.name == "Success - Register citizen" {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), result)
				assert.Empty(suite.T(), result.Password)
			} else {
				assert.NotNil(suite.T(), err)
				assert.Nil(suite.T(), result)
			}
		})
	}
}

func (suite *AuthServiceTestSuite) TestLogin() {
	input := service.LoginInput{
		Email:    "rian@mail.com",
		Password: "password123",
	}

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		suite.T().Fatal("failed to generate dummy hash for testing:", err)
	}

	tests := []struct {
		name   string
		mockFn func()
	}{
		{
			name: "Success - Login citizen",
			mockFn: func() {
				userData := &entity.User{
					ID:       10,
					Email:    input.Email,
					Password: hashedPassword,
					Role:     entity.RoleCitizen,
				}
				suite.mockUserRepo.EXPECT().
					FindByEmail(input.Email).
					Return(userData, nil)
			},
		},
		{
			name: "Fail - User email not found",
			mockFn: func() {
				suite.mockUserRepo.EXPECT().
					FindByEmail(input.Email).
					Return(nil, errors.New("user not found"))
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockFn()

			svc := service.NewAuthService(suite.mockUserRepo, suite.cfg, suite.mailer)
			token, err := svc.Login(input)

			if tc.name == "Success - Login citizen" {
				assert.Nil(suite.T(), err)
				assert.NotEmpty(suite.T(), token)
			} else {
				assert.NotNil(suite.T(), err)
				assert.Empty(suite.T(), token)
			}
		})
	}
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
