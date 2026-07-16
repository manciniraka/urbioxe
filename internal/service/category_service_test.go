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

type CategoryServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockCategory *mocks.MockCategoryRepository
}

func (suite *CategoryServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCategory = mocks.NewMockCategoryRepository(suite.ctrl)
}

func (suite *CategoryServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *CategoryServiceTestSuite) TestGetAllCategories() {
	suite.Run("Success - Get all active categories", func() {
		suite.SetupTest()
		expectedData := []entity.Category{
			{ID: 1, Name: "Jalan Rusak", IsActive: true},
			{ID: 2, Name: "Sampah", IsActive: true},
		}

		suite.mockCategory.EXPECT().FindAll(true).Return(expectedData, nil)

		svc := service.NewCategoryService(suite.mockCategory)
		res, err := svc.GetAllCategories(true)

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 2)
	})
}

func (suite *CategoryServiceTestSuite) TestGetCategoryByID() {
	suite.Run("Success - Found", func() {
		suite.SetupTest()
		category := &entity.Category{ID: 1, Name: "Banjir"}

		suite.mockCategory.EXPECT().FindByID(uint(1)).Return(category, nil)

		svc := service.NewCategoryService(suite.mockCategory)
		res, err := svc.GetCategoryByID(1)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Banjir", res.Name)
	})

	suite.Run("Fail - Not Found", func() {
		suite.SetupTest()
		suite.mockCategory.EXPECT().FindByID(uint(99)).Return(nil, gorm.ErrRecordNotFound)

		svc := service.NewCategoryService(suite.mockCategory)
		res, err := svc.GetCategoryByID(99)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *CategoryServiceTestSuite) TestCreateCategory() {
	input := service.CreateCategoryInput{
		DepartmentID: 2,
		Name:         "Lampu Jalan Mati",
		Description:  "Fasilitas penerangan umum padam",
	}

	suite.Run("Success - Created by admin", func() {
		suite.SetupTest()
		suite.mockCategory.EXPECT().
			CheckNameExistsInDepartment(input.DepartmentID, input.Name, uint(0)).
			Return(false, nil)

		suite.mockCategory.EXPECT().Create(gomock.Any()).Return(nil)

		svc := service.NewCategoryService(suite.mockCategory)
		res, err := svc.CreateCategory("super_admin", input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.True(suite.T(), res.IsActive)
	})

	suite.Run("Fail - Forbidden Role", func() {
		suite.SetupTest()
		svc := service.NewCategoryService(suite.mockCategory)
		res, err := svc.CreateCategory("citizen", input)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})

	suite.Run("Fail - Name Already Exists", func() {
		suite.SetupTest()
		suite.mockCategory.EXPECT().
			CheckNameExistsInDepartment(input.DepartmentID, input.Name, uint(0)).
			Return(true, nil)

		svc := service.NewCategoryService(suite.mockCategory)
		res, err := svc.CreateCategory("department_admin", input)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *CategoryServiceTestSuite) TestUpdateCategory() {
	categoryID := uint(5)
	input := service.UpdateCategoryInput{
		DepartmentID: 3,
		Name:         "Pohon Tumbang",
		Description:  "Menghalangi jalan raya",
	}

	suite.Run("Success - Updated", func() {
		suite.SetupTest()
		existingData := &entity.Category{ID: categoryID, Name: "Pohon Rusak", DepartmentID: 1}

		suite.mockCategory.EXPECT().FindByID(categoryID).Return(existingData, nil)

		suite.mockCategory.EXPECT().
			CheckNameExistsInDepartment(input.DepartmentID, input.Name, categoryID).
			Return(false, nil)

		suite.mockCategory.EXPECT().Update(gomock.Any()).Return(nil)

		svc := service.NewCategoryService(suite.mockCategory)
		res, err := svc.UpdateCategory(categoryID, "super_admin", input)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Pohon Tumbang", res.Name)
	})
}

func (suite *CategoryServiceTestSuite) TestToggleCategoryStatus() {
	categoryID := uint(10)
	input := service.ToggleCategoryStatusInput{IsActive: false}

	suite.Run("Success - Deactivate Category", func() {
		suite.SetupTest()
		suite.mockCategory.EXPECT().FindByID(categoryID).Return(&entity.Category{ID: categoryID}, nil)
		suite.mockCategory.EXPECT().UpdateStatus(categoryID, false).Return(nil)

		svc := service.NewCategoryService(suite.mockCategory)
		err := svc.ToggleCategoryStatus(categoryID, "department_admin", input)

		assert.Nil(suite.T(), err)
	})
}

func TestCategoryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceTestSuite))
}
