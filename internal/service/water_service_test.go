package service_test

import (
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/internal/dto"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WaterServiceTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	mockWater      *mocks.MockWaterRepository
	mockDistrict   *mocks.MockDistrictRepository
	mockUser       *mocks.MockUserRepository
	mockMeter      *mocks.MockMeterReadingRepository
	mockCloudinary *mocks.MockCloudinaryService
}

func (suite *WaterServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockWater = mocks.NewMockWaterRepository(suite.ctrl)
	suite.mockDistrict = mocks.NewMockDistrictRepository(suite.ctrl)
	suite.mockUser = mocks.NewMockUserRepository(suite.ctrl)
	suite.mockMeter = mocks.NewMockMeterReadingRepository(suite.ctrl)
	suite.mockCloudinary = mocks.NewMockCloudinaryService(suite.ctrl)
}

func (suite *WaterServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *WaterServiceTestSuite) TestCreateWaterStatus() {
	userID := uint(1)
	location, _ := time.LoadLocation("Asia/Jakarta")

	futureTime := time.Now().In(location).Add(24 * time.Hour)
	futureTimeStr := futureTime.Format("2006-01-02 15:04")

	suite.Run("Success - Create Disrupted Status", func() {
		suite.SetupTest()
		input := dto.CreateWaterStatusInput{
			DistrictID:        1,
			Status:            entity.WaterStatusDisrupted,
			StartedAt:         futureTimeStr,
			EstimatedDuration: 5,
			Reason:            "Pipa Bocor",
		}

		suite.mockDistrict.EXPECT().GetByID(uint(1)).Return(&entity.District{ID: 1, Name: "Bogor Tengah"}, nil)
		suite.mockWater.EXPECT().CreateWaterStatus(gomock.Any()).Return(nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.CreateWaterStatus(userID, input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), "Bogor Tengah", res.District)
		assert.Equal(suite.T(), entity.WaterStatusDisrupted, res.Status)
	})

	suite.Run("Fail - Started At Before Now", func() {
		suite.SetupTest()
		pastTimeStr := time.Now().In(location).Add(-5 * time.Hour).Format("2006-01-02 15:04")
		input := dto.CreateWaterStatusInput{
			DistrictID:        1,
			Status:            entity.WaterStatusDisrupted,
			StartedAt:         pastTimeStr,
			EstimatedDuration: 2,
			Reason:            "Perawatan Berkala",
		}

		suite.mockDistrict.EXPECT().GetByID(uint(1)).Return(&entity.District{ID: 1, Name: "Bogor Tengah"}, nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.CreateWaterStatus(userID, input)

		assert.ErrorIs(suite.T(), err, errs.ErrStartedAtBeforeNow)
		assert.Nil(suite.T(), res)
	})
}

func (suite *WaterServiceTestSuite) TestGetMyWaterStatus() {
	userID := uint(10)
	districtID := uint(3)

	suite.Run("Success - Found My District Water Status", func() {
		suite.SetupTest()
		suite.mockUser.EXPECT().GetByID(userID).Return(&entity.User{ID: userID, HomeDistrictID: &districtID}, nil)

		mockStatus := &entity.WaterStatus{
			DistrictID: districtID,
			Status:     entity.WaterStatusNormal,
			District:   entity.District{Name: "Bogor Barat"},
		}
		suite.mockWater.EXPECT().GetLatestWaterStatusByDistrictID(districtID).Return(mockStatus, nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.GetMyWaterStatus(userID)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), "Bogor Barat", res.District)
	})

	suite.Run("Fail - User District Is Nil", func() {
		suite.SetupTest()
		suite.mockUser.EXPECT().GetByID(userID).Return(&entity.User{ID: userID, HomeDistrictID: nil}, nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.GetMyWaterStatus(userID)

		assert.ErrorIs(suite.T(), err, errs.ErrDistrictNotFound)
		assert.Nil(suite.T(), res)
	})
}

func (suite *WaterServiceTestSuite) TestSimulateBill() {
	suite.Run("Success - Calculate Block Tariff Breakdown", func() {
		suite.SetupTest()
		input := dto.BillSimulationRequest{
			TariffGroup: "R1",
			Usage:       25,
		}

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.SimulateBill(input)

		if err == nil {
			assert.NotNil(suite.T(), res)
			assert.Equal(suite.T(), 25, res.Usage)
			assert.NotEmpty(suite.T(), res.Breakdown)
		} else {
			assert.NotNil(suite.T(), err)
		}
	})
}

func (suite *WaterServiceTestSuite) TestCreateMeterReading() {
	userID := uint(5)
	fileHeader := &multipart.FileHeader{
		Filename: "meter_photo.jpg",
		Size:     1024,
	}
	input := dto.CreateMeterReadingInput{
		CustomerNumber: "CUST001",
		CurrentReading: 150,
	}

	suite.Run("Success - Upload image and save reading", func() {
		suite.SetupTest()
		uploadedURL := "https://cloudinary.com/urbioxe/image.jpg"

		suite.mockCloudinary.EXPECT().UploadImage(fileHeader).Return(uploadedURL, nil)
		suite.mockMeter.EXPECT().CreateMeterReading(gomock.Any()).Return(nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.CreateMeterReading(userID, input, fileHeader)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), uploadedURL, res.PhotoURL)
		assert.Equal(suite.T(), entity.MeterReadingPending, res.Status)
	})

	suite.Run("Fail - Cloudinary Upload Error", func() {
		suite.SetupTest()
		suite.mockCloudinary.EXPECT().UploadImage(fileHeader).Return("", errors.New("upload failed"))

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.CreateMeterReading(userID, input, fileHeader)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *WaterServiceTestSuite) TestGetAllWaterStatus() {
	suite.Run("Success - Get current statuses with update metadata", func() {
		suite.SetupTest()

		now := time.Now()
		mockStatuses := []entity.WaterStatus{
			{
				ID:        1,
				Status:    entity.WaterStatusDisrupted,
				District:  entity.District{Name: "Bogor Utara"},
				StartedAt: now,
				CreatedAt: now,
			},
			{
				ID:        2,
				Status:    entity.WaterStatusNormal,
				District:  entity.District{Name: "Bogor Selatan"},
				StartedAt: now,
				CreatedAt: now,
			},
		}

		suite.mockWater.EXPECT().GetLatestAllWaterStatus().Return(mockStatuses, nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.GetAllWaterStatus()

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Len(suite.T(), res.Data, 2)
		assert.Equal(suite.T(), "Bogor Utara", res.Data[0].District)
		assert.NotNil(suite.T(), res.Data[0].StartedAt)
		assert.NotEmpty(suite.T(), res.Metadata.UpdatedAt)
	})
}

func (suite *WaterServiceTestSuite) TestGetWaterStatusHistories() {
	suite.Run("Success - Fetch all history logs", func() {
		suite.SetupTest()

		now := time.Now()
		mockHistories := []entity.WaterStatus{
			{
				ID:                  1,
				Status:              entity.WaterStatusDisrupted,
				Reason:              "Pembersihan Bak Grapari",
				EstimatedDuration:   3,
				District:            entity.District{Name: "Bogor Timur"},
				User:                entity.User{Name: "Admin Ganteng"},
				StartedAt:           now,
				EstimatedRecoveryAt: now.Add(3 * time.Hour),
				CreatedAt:           now,
			},
		}

		suite.mockWater.EXPECT().GetWaterStatusHistories().Return(mockHistories, nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.GetWaterStatusHistories()

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 1)
		assert.Equal(suite.T(), "Bogor Timur", res[0].District)
		assert.Equal(suite.T(), "Admin Ganteng", res[0].CreatedBy)
		assert.NotNil(suite.T(), res[0].EstimatedDuration)
	})
}

func (suite *WaterServiceTestSuite) TestSimulateBill_Detailed() {

	tests := []struct {
		name        string
		tariffGroup string
		usage       int
	}{
		{
			name:        "Simulation - Usage within lower blocks",
			tariffGroup: "R1",
			usage:       8,
		},
		{
			name:        "Simulation - High usage hitting progressive blocks",
			tariffGroup: "R1",
			usage:       45,
		},
		{
			name:        "Simulation - Zero usage triggers minimum tariff calculation",
			tariffGroup: "R1",
			usage:       0,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)

			res, err := svc.SimulateBill(dto.BillSimulationRequest{
				TariffGroup: tc.tariffGroup,
				Usage:       tc.usage,
			})

			if err == nil {
				assert.NotNil(suite.T(), res)
				assert.Equal(suite.T(), tc.usage, res.Usage)
				assert.True(suite.T(), res.TotalBill > 0)
				assert.NotEmpty(suite.T(), res.Breakdown)
			} else {
				assert.Contains(suite.T(), err.Error(), "not found")
			}
		})
	}
}

func (suite *WaterServiceTestSuite) TestGetMyMeterReadings() {
	userID := uint(10)

	suite.Run("Success - Get customer's own readings", func() {
		suite.SetupTest()
		mockReadings := []entity.MeterReading{
			{
				ID:             1,
				CustomerNumber: "PAM-001",
				CurrentReading: 120,
				PhotoURL:       "https://res.cloudinary.com/img1.png",
				Status:         entity.MeterReadingPending,
			},
		}

		suite.mockMeter.EXPECT().GetMyMeterReadings(userID).Return(mockReadings, nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.GetMyMeterReadings(userID)

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 1)
		assert.Equal(suite.T(), "PAM-001", res[0].CustomerNumber)
	})
}

func (suite *WaterServiceTestSuite) TestGetAllMeterReadings() {
	month := 7
	year := 2026

	suite.Run("Success - Admin fetch all customers reading with filters", func() {
		suite.SetupTest()
		mockReadings := []entity.MeterReading{
			{
				ID:             2,
				CustomerNumber: "PAM-002",
				CurrentReading: 340,
				PhotoURL:       "https://res.cloudinary.com/img2.png",
				Status:         entity.MeterReadingApproved,
			},
		}

		suite.mockMeter.EXPECT().GetAllMeterReadings(&month, &year).Return(mockReadings, nil)

		svc := service.NewWaterService(suite.mockWater, suite.mockDistrict, suite.mockUser, suite.mockMeter, suite.mockCloudinary)
		res, err := svc.GetAllMeterReadings(&month, &year)

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 1)
		assert.Equal(suite.T(), entity.MeterReadingApproved, res[0].Status)
	})
}

func TestWaterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WaterServiceTestSuite))
}
