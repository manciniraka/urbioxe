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

type EmergencyServiceTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	mockEmergency *mocks.MockEmergencyRepository
}

func (suite *EmergencyServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockEmergency = mocks.NewMockEmergencyRepository(suite.ctrl)
}

func (suite *EmergencyServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func stringPtr(s string) *string {
	return &s
}

func (suite *EmergencyServiceTestSuite) TestGetAll() {
	suite.Run("Success - Get all emergency contacts", func() {
		suite.SetupTest()
		expectedData := []entity.EmergencyContact{
			{ID: 1, Name: "Ambulan RSUD", PhoneNumber: "118"},
			{ID: 2, Name: "Pemadam Kebakaran", PhoneNumber: "113"},
		}

		suite.mockEmergency.EXPECT().GetAll().Return(expectedData, nil)

		svc := service.NewEmergencyService(suite.mockEmergency)
		res, err := svc.GetAll()

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 2)
	})
}

func (suite *EmergencyServiceTestSuite) TestGetByID() {
	suite.Run("Success - Found", func() {
		suite.SetupTest()
		contact := &entity.EmergencyContact{ID: 1, Name: "Polsek Bogor"}

		suite.mockEmergency.EXPECT().GetByID(uint(1)).Return(contact, nil)

		svc := service.NewEmergencyService(suite.mockEmergency)
		res, err := svc.GetByID(1)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Polsek Bogor", res.Name)
	})

	suite.Run("Fail - Not Found", func() {
		suite.SetupTest()
		suite.mockEmergency.EXPECT().GetByID(uint(99)).Return(nil, gorm.ErrRecordNotFound)

		svc := service.NewEmergencyService(suite.mockEmergency)
		res, err := svc.GetByID(99)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *EmergencyServiceTestSuite) TestCreate() {
	deptID := uint(2)
	distID := uint(4)
	input := service.CreateEmergencyInput{
		DepartmentID: &deptID,
		DistrictID:   &distID,
		Name:         "Damkar Sektor Barat",
		PhoneNumber:  "0251-123456",
		Description:  stringPtr("Kontak darurat kebakaran barat"),
		IconURL:      stringPtr("https://cdn.com/icon.png"),
		IsActive:     true,
	}

	suite.Run("Success - Created emergency contact", func() {
		suite.SetupTest()
		suite.mockEmergency.EXPECT().Create(gomock.Any()).Return(nil)

		svc := service.NewEmergencyService(suite.mockEmergency)
		res, err := svc.Create(input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), "Damkar Sektor Barat", res.Name)
	})
}

func (suite *EmergencyServiceTestSuite) TestUpdate() {
	contactID := uint(5)
	deptID := uint(2)
	distID := uint(4)
	input := service.UpdateEmergencyInput{
		DepartmentID: &deptID,
		DistrictID:   &distID,
		Name:         "Damkar Sektor Barat Diperbarui",
		PhoneNumber:  "0251-654321",
		Description:  stringPtr("Kontak darurat kebakaran barat revisi"),
		IconURL:      stringPtr("https://cdn.com/icon-new.png"),
		IsActive:     false,
	}

	suite.Run("Success - Updated fully", func() {
		suite.SetupTest()
		existingData := &entity.EmergencyContact{ID: contactID, Name: "Damkar Sektor Barat Lama"}

		suite.mockEmergency.EXPECT().GetByID(contactID).Return(existingData, nil)
		suite.mockEmergency.EXPECT().Update(gomock.Any()).Return(nil)

		svc := service.NewEmergencyService(suite.mockEmergency)
		res, err := svc.Update(contactID, input)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Damkar Sektor Barat Diperbarui", res.Name)
		assert.False(suite.T(), res.IsActive)
	})

	suite.Run("Fail - Update Not Found", func() {
		suite.SetupTest()
		suite.mockEmergency.EXPECT().GetByID(contactID).Return(nil, gorm.ErrRecordNotFound)

		svc := service.NewEmergencyService(suite.mockEmergency)
		res, err := svc.Update(contactID, input)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func TestEmergencyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EmergencyServiceTestSuite))
}
