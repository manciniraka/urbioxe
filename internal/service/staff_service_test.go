package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StaffServiceTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	mockUserRepo  *mocks.MockUserRepository
	mockStaffRepo *mocks.MockStaffRepository
	mailer        *mailjet.Client
	mockServer    *httptest.Server
	gormDB        *gorm.DB
	sqlMock       sqlmock.Sqlmock
}

func (suite *StaffServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockUserRepo = mocks.NewMockUserRepository(suite.ctrl)
	suite.mockStaffRepo = mocks.NewMockStaffRepository(suite.ctrl)

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

	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("failed to open sqlmock: %s", err)
	}
	suite.sqlMock = mock

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("failed to open gorm: %s", err)
	}
	suite.gormDB = gormDB
}

func (suite *StaffServiceTestSuite) TearDownTest() {
	suite.mockServer.Close()
	suite.ctrl.Finish()
}

func (suite *StaffServiceTestSuite) TestCreateStaff() {
	input := service.CreateStaffInput{
		Position:     entity.PositionFieldOfficer,
		Email:        "officer@urbioxe.com",
		NIK:          "3201010101010005",
		Name:         "Budi Officer",
		PhoneNumber:  "081234567890",
		DepartmentID: 1,
		JoinDate:     "2026-07-16",
	}

	suite.Run("Success - Create Staff", func() {
		suite.SetupTest()

		suite.mockUserRepo.EXPECT().FindByEmail(input.Email).Return(nil, gorm.ErrRecordNotFound)
		suite.mockUserRepo.EXPECT().FindByNIK(input.NIK).Return(nil, gorm.ErrRecordNotFound)

		suite.sqlMock.ExpectQuery(`SELECT \* FROM "departments"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).AddRow(1, "DISHUB", "Dinas Perhubungan"))

		suite.mockStaffRepo.EXPECT().
			GetLastEmployeeSequence(gomock.Any(), gomock.Any()).
			Return(5, nil)

		suite.sqlMock.ExpectBegin()
		suite.mockUserRepo.EXPECT().RegisterUserTx(gomock.Any(), gomock.Any()).DoAndReturn(func(tx *gorm.DB, u *entity.User) error {
			u.ID = 10
			return nil
		})
		suite.mockStaffRepo.EXPECT().CreateStaffTx(gomock.Any(), gomock.Any()).Return(nil)
		suite.sqlMock.ExpectCommit()

		svc := service.NewStaffService(suite.gormDB, suite.mockStaffRepo, suite.mockUserRepo, suite.mailer)
		res, err := svc.CreateStaff(input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), uint(10), res.UserID)
		assert.Equal(suite.T(), "Budi Officer", res.User.Name)
	})

	suite.Run("Fail - Invalid Position", func() {
		suite.SetupTest()
		invalidInput := input
		invalidInput.Position = "INVALID_POSITION"

		svc := service.NewStaffService(suite.gormDB, suite.mockStaffRepo, suite.mockUserRepo, suite.mailer)
		res, err := svc.CreateStaff(invalidInput)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *StaffServiceTestSuite) TestGetAllStaff() {
	suite.Run("Success - Get All", func() {
		suite.SetupTest()
		expectedData := []entity.StaffProfile{
			{
				ID:             1,
				EmployeeNumber: "DISHUB-2026-001",
				DepartmentID:   1,
				Position:       entity.PositionFieldOfficer,
				IsActive:       true,
				User:           &entity.User{Name: "Budi"},
				Department:     &entity.Department{Name: "Dinas Perhubungan"},
			},
		}

		suite.mockStaffRepo.EXPECT().GetAll().Return(expectedData, nil)

		svc := service.NewStaffService(suite.gormDB, suite.mockStaffRepo, suite.mockUserRepo, suite.mailer)
		res, err := svc.GetAllStaff()

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 1)
		assert.Equal(suite.T(), "Budi", res[0].Name)
	})
}

func (suite *StaffServiceTestSuite) TestGetStaffByID() {
	suite.Run("Success - Found", func() {
		suite.SetupTest()
		expectedData := &entity.StaffProfile{
			ID:             1,
			UserID:         10,
			DepartmentID:   1,
			EmployeeNumber: "DISHUB-2026-001",
			Position:       entity.PositionFieldOfficer,
			JoinDate:       "2026-07-16",
			IsActive:       true,
			User:           &entity.User{NIK: "123", Name: "Budi", Email: "budi@mail.com", PhoneNumber: "081"},
			Department:     &entity.Department{Name: "Dinas Perhubungan"},
		}

		suite.mockStaffRepo.EXPECT().GetByID(uint(1)).Return(expectedData, nil)

		svc := service.NewStaffService(suite.gormDB, suite.mockStaffRepo, suite.mockUserRepo, suite.mailer)
		res, err := svc.GetStaffByID(1)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), "Budi", res.Name)
	})
}

func (suite *StaffServiceTestSuite) TestUpdateStaff() {
	input := service.UpdateStaffInput{
		DepartmentID: 1,
		Position:     entity.PositionDepartmentAdmin,
		PhoneNumber:  "0899999",
		IsActive:     true,
	}

	suite.Run("Success - Update Staff", func() {
		suite.SetupTest()
		existingStaff := &entity.StaffProfile{
			ID:           1,
			UserID:       10,
			DepartmentID: 1,
			Position:     entity.PositionFieldOfficer,
			User:         &entity.User{ID: 10, PhoneNumber: "081"},
		}

		suite.mockStaffRepo.EXPECT().GetByID(uint(1)).Return(existingStaff, nil)

		suite.sqlMock.ExpectQuery(`SELECT \* FROM "departments"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name"}).AddRow(1, "DISHUB", "Dinas Perhubungan"))

		suite.sqlMock.ExpectBegin()
		suite.mockUserRepo.EXPECT().UpdateUserTx(gomock.Any(), gomock.Any()).Return(nil)
		suite.mockStaffRepo.EXPECT().UpdateStaffTx(gomock.Any(), gomock.Any()).Return(nil)
		suite.sqlMock.ExpectCommit()

		suite.mockStaffRepo.EXPECT().GetByID(uint(1)).Return(existingStaff, nil)

		svc := service.NewStaffService(suite.gormDB, suite.mockStaffRepo, suite.mockUserRepo, suite.mailer)
		res, err := svc.UpdateStaff(1, input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
	})
}

func TestStaffServiceTestSuite(t *testing.T) {
	suite.Run(t, new(StaffServiceTestSuite))
}
