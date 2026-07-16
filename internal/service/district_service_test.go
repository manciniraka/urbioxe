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
)

type DistrictServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockDistrict *mocks.MockDistrictRepository
}

func (suite *DistrictServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDistrict = mocks.NewMockDistrictRepository(suite.ctrl)
}

func (suite *DistrictServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *DistrictServiceTestSuite) TestGetAllDistrict() {
	suite.Run("Success - Get all districts", func() {
		suite.SetupTest()
		expectedData := []entity.District{
			{ID: 1, Name: "Kecamatan Bogor Tengah"},
			{ID: 2, Name: "Kecamatan Bogor Barat"},
		}

		suite.mockDistrict.EXPECT().GetAll().Return(expectedData, nil)

		svc := service.NewDistrictService(suite.mockDistrict)
		res, err := svc.GetAllDistrict()

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 2)
	})
}

func (suite *DistrictServiceTestSuite) TestGetDistrictByID() {
	suite.Run("Success - Found", func() {
		suite.SetupTest()
		district := &entity.District{ID: 1, Name: "Kecamatan Bogor Tengah"}

		suite.mockDistrict.EXPECT().GetByID(uint(1)).Return(district, nil)

		svc := service.NewDistrictService(suite.mockDistrict)
		res, err := svc.GetDistrictByID(1)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Kecamatan Bogor Tengah", res.Name)
	})
}

func (suite *DistrictServiceTestSuite) TestCreateDistrict() {
	input := service.CreateDistrictInput{
		Name: "Kecamatan Bogor Timur",
	}

	suite.Run("Success - Created new district", func() {
		suite.SetupTest()

		suite.mockDistrict.EXPECT().
			GetByName(input.Name).
			Return(nil, errs.ErrDistrictNotFound)

		suite.mockDistrict.EXPECT().Create(gomock.Any()).Return(nil)

		svc := service.NewDistrictService(suite.mockDistrict)
		res, err := svc.CreateDistrict(input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), "Kecamatan Bogor Timur", res.Name)
	})

	suite.Run("Fail - District Already Exists", func() {
		suite.SetupTest()
		existingData := &entity.District{ID: 1, Name: input.Name}
		suite.mockDistrict.EXPECT().
			GetByName(input.Name).
			Return(existingData, nil)

		svc := service.NewDistrictService(suite.mockDistrict)
		res, err := svc.CreateDistrict(input)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *DistrictServiceTestSuite) TestUpdateDistrict() {
	districtID := uint(5)
	input := service.UpdateDistrictInput{
		IsActive: false,
	}

	suite.Run("Success - Updated status", func() {
		suite.SetupTest()
		existingData := &entity.District{ID: districtID, Name: "Kecamatan Bogor Selatan", IsActive: true}

		suite.mockDistrict.EXPECT().GetByID(districtID).Return(existingData, nil)
		suite.mockDistrict.EXPECT().Update(gomock.Any()).Return(nil)

		svc := service.NewDistrictService(suite.mockDistrict)
		res, err := svc.UpdateDistrict(districtID, input)

		assert.Nil(suite.T(), err)
		assert.False(suite.T(), res.IsActive)
	})
}

func TestDistrictServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DistrictServiceTestSuite))
}
