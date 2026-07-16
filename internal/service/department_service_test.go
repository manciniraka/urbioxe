package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type DepartmentServiceTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	mockDepartment *mocks.MockDepartmentRepository
}

func (suite *DepartmentServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDepartment = mocks.NewMockDepartmentRepository(suite.ctrl)
}

func (suite *DepartmentServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *DepartmentServiceTestSuite) TestGetAllDepartments() {
	suite.Run("Success - Get all active departments", func() {
		suite.SetupTest()
		expectedData := []entity.Department{
			{ID: 1, Code: "DISHUB", Name: "Dinas Perhubungan", IsActive: true},
			{ID: 2, Code: "DLH", Name: "Dinas Lingkungan Hidup", IsActive: true},
		}

		suite.mockDepartment.EXPECT().FindAll(true).Return(expectedData, nil)

		svc := service.NewDepartmentService(suite.mockDepartment)
		res, err := svc.GetAllDepartments(true)

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 2)
	})
}

func (suite *DepartmentServiceTestSuite) TestGetDepartmentByID() {
	suite.Run("Success - Found", func() {
		suite.SetupTest()
		dept := &entity.Department{ID: 1, Code: "DISHUB", Name: "Dinas Perhubungan"}

		suite.mockDepartment.EXPECT().FindByID(uint(1)).Return(dept, nil)

		svc := service.NewDepartmentService(suite.mockDepartment)
		res, err := svc.GetDepartmentByID(1)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Dinas Perhubungan", res.Name)
	})

	suite.Run("Fail - Not Found", func() {
		suite.SetupTest()
		suite.mockDepartment.EXPECT().FindByID(uint(99)).Return(nil, gorm.ErrRecordNotFound)

		svc := service.NewDepartmentService(suite.mockDepartment)
		res, err := svc.GetDepartmentByID(99)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *DepartmentServiceTestSuite) TestCreateDepartment() {
	input := service.DepartmentInput{
		Code:        "DISHUB",
		Name:        "Dinas Perhubungan",
		Description: "Mengurus transportasi publik",
	}

	suite.Run("Success - Created by super_admin", func() {
		suite.SetupTest()
		suite.mockDepartment.EXPECT().
			CheckCodeOrNameExists(input.Code, input.Name, uint(0)).
			Return(false, nil)

		suite.mockDepartment.EXPECT().Create(gomock.Any()).Return(nil)

		svc := service.NewDepartmentService(suite.mockDepartment)
		res, err := svc.CreateDepartment("super_admin", input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.True(suite.T(), res.IsActive)
	})

	suite.Run("Fail - Forbidden Role", func() {
		suite.SetupTest()
		svc := service.NewDepartmentService(suite.mockDepartment)
		res, err := svc.CreateDepartment("department_admin", input)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})

	suite.Run("Fail - Code or Name Already Used", func() {
		suite.SetupTest()
		suite.mockDepartment.EXPECT().
			CheckCodeOrNameExists(input.Code, input.Name, uint(0)).
			Return(true, nil)

		svc := service.NewDepartmentService(suite.mockDepartment)
		res, err := svc.CreateDepartment("super_admin", input)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *DepartmentServiceTestSuite) TestUpdateDepartment() {
	deptID := uint(5)
	input := service.DepartmentInput{
		Code:        "DLH",
		Name:        "Dinas Lingkungan Hidup Baru",
		Description: "Mengurus kebersihan wilayah",
	}

	suite.Run("Success - Updated", func() {
		suite.SetupTest()
		existingData := &entity.Department{ID: deptID, Code: "DLH", Name: "Dinas Lingkungan Hidup"}

		suite.mockDepartment.EXPECT().FindByID(deptID).Return(existingData, nil)
		suite.mockDepartment.EXPECT().
			CheckCodeOrNameExists(input.Code, input.Name, deptID).
			Return(false, nil)

		suite.mockDepartment.EXPECT().Update(gomock.Any()).Return(nil)

		svc := service.NewDepartmentService(suite.mockDepartment)
		res, err := svc.UpdateDepartment(deptID, "super_admin", input)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Dinas Lingkungan Hidup Baru", res.Name)
	})
}

func (suite *DepartmentServiceTestSuite) TestToggleDepartmentStatus() {
	deptID := uint(10)
	input := service.ToggleDepartmentStatusInput{IsActive: false}

	suite.Run("Success - Deactivate Department", func() {
		suite.SetupTest()
		suite.mockDepartment.EXPECT().FindByID(deptID).Return(&entity.Department{ID: deptID}, nil)
		suite.mockDepartment.EXPECT().UpdateStatus(deptID, false).Return(nil)

		svc := service.NewDepartmentService(suite.mockDepartment)
		err := svc.ToggleDepartmentStatus(deptID, "super_admin", input)

		assert.Nil(suite.T(), err)
	})
}

func TestDepartmentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DepartmentServiceTestSuite))
}
