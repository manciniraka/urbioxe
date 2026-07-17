package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
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

	tests := []struct {
		name     string
		role     string
		mockFn   func()
		verifyFn func(res *entity.Category, err error)
	}{
		{
			name: "Success - Updated fully and refetched",
			role: "super_admin",
			mockFn: func() {
				existingData := &entity.Category{ID: categoryID, Name: "Pohon Rusak", DepartmentID: 1}
				suite.mockCategory.EXPECT().FindByID(categoryID).Return(existingData, nil)

				suite.mockCategory.EXPECT().
					CheckNameExistsInDepartment(input.DepartmentID, input.Name, categoryID).
					Return(false, nil)

				suite.mockCategory.EXPECT().Update(gomock.Any()).Return(nil)

				updatedData := &entity.Category{
					ID:           categoryID,
					Name:         input.Name,
					Description:  input.Description,
					DepartmentID: input.DepartmentID,
				}
				suite.mockCategory.EXPECT().FindByID(categoryID).Return(updatedData, nil)
			},
			verifyFn: func(res *entity.Category, err error) {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), res)
				assert.Equal(suite.T(), "Pohon Tumbang", res.Name)
				assert.Equal(suite.T(), uint(3), res.DepartmentID)
			},
		},
		{
			name:   "Fail - Forbidden Role",
			role:   "citizen",
			mockFn: func() {},
			verifyFn: func(res *entity.Category, err error) {
				assert.ErrorIs(suite.T(), err, errs.ErrForbidden)
				assert.Nil(suite.T(), res)
			},
		},
		{
			name: "Fail - Category Not Found",
			role: "department_admin",
			mockFn: func() {
				suite.mockCategory.EXPECT().FindByID(categoryID).Return(nil, gorm.ErrRecordNotFound)
			},
			verifyFn: func(res *entity.Category, err error) {
				assert.ErrorIs(suite.T(), err, errs.ErrCategoryNotFound)
				assert.Nil(suite.T(), res)
			},
		},
		{
			name: "Fail - Name Already Exists In Department",
			role: "department_admin",
			mockFn: func() {
				existingData := &entity.Category{ID: categoryID, Name: "Pohon Rusak", DepartmentID: 1}
				suite.mockCategory.EXPECT().FindByID(categoryID).Return(existingData, nil)

				suite.mockCategory.EXPECT().
					CheckNameExistsInDepartment(input.DepartmentID, input.Name, categoryID).
					Return(true, nil)
			},
			verifyFn: func(res *entity.Category, err error) {
				assert.ErrorIs(suite.T(), err, errs.ErrCategoryAlreadyExists)
				assert.Nil(suite.T(), res)
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockFn()

			svc := service.NewCategoryService(suite.mockCategory)
			res, err := svc.UpdateCategory(categoryID, tc.role, input)

			tc.verifyFn(res, err)
		})
	}
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
