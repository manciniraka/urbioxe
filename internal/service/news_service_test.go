package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func int64Ptr(n int64) *int64 {
	return &n
}

type NewsServiceTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	mockNews *mocks.MockNewsRepository
}

func (suite *NewsServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockNews = mocks.NewMockNewsRepository(suite.ctrl)
}

func (suite *NewsServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *NewsServiceTestSuite) TestGetAll() {
	suite.Run("Success - Get all news", func() {
		suite.SetupTest()
		expectedData := []entity.RegionalNews{
			{ID: 1, Title: "Pesta Rakyat Bogor"},
			{ID: 2, Title: "Perbaikan Jalan Juanda"},
		}

		suite.mockNews.EXPECT().GetAll().Return(expectedData, nil)

		svc := service.NewNewsService(suite.mockNews)
		res, err := svc.GetAll()

		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), res, 2)
	})
}

func (suite *NewsServiceTestSuite) TestGetByID() {
	suite.Run("Success - Found", func() {
		suite.SetupTest()
		news := &entity.RegionalNews{ID: 1, Title: "Pesta Rakyat Bogor"}

		suite.mockNews.EXPECT().GetByID(int64(1)).Return(news, nil)

		svc := service.NewNewsService(suite.mockNews)
		res, err := svc.GetByID(1)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Pesta Rakyat Bogor", res.Title)
	})
}

func (suite *NewsServiceTestSuite) TestCreate() {
	createdBy := int64(10)

	input := service.CreateNewsInput{
		DepartmentID: int64Ptr(2),
		DistrictID:   int64Ptr(4),
		Title:        "Info Pemadaman Air",
		Content:      "Pemadaman air berkala di sektor utara mulai malam ini.",
		Category:     "announcement",
		BannerURL:    stringPtr("https://cdn.com/banner.png"),
		TargetScope:  "district",
		IsPinned:     false,
	}

	suite.Run("Success - Created news", func() {
		suite.SetupTest()
		suite.mockNews.EXPECT().Create(gomock.Any()).Return(nil)

		svc := service.NewNewsService(suite.mockNews)
		res, err := svc.Create(createdBy, input)

		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), res)
		assert.Equal(suite.T(), "Info Pemadaman Air", res.Title)
		assert.Equal(suite.T(), &createdBy, res.CreatedBy)
	})
}

func (suite *NewsServiceTestSuite) TestUpdate() {
	newsID := int64(5)

	input := service.UpdateNewsInput{
		DepartmentID: int64Ptr(2),
		DistrictID:   int64Ptr(4),
		Title:        "Info Pemadaman Air Diperbarui",
		Content:      "Pemadaman air berkala sektor utara dibatalkan.",
		Category:     "Pengumuman",
		BannerURL:    stringPtr("https://cdn.com/banner-rev.png"),
		TargetScope:  "district",
		IsPinned:     true,
	}

	suite.Run("Success - Updated successfully", func() {
		suite.SetupTest()
		existingData := &entity.RegionalNews{ID: newsID, Title: "Info Pemadaman Air Lama"}

		suite.mockNews.EXPECT().GetByID(newsID).Return(existingData, nil)
		suite.mockNews.EXPECT().Update(gomock.Any()).Return(nil)

		svc := service.NewNewsService(suite.mockNews)
		res, err := svc.Update(newsID, input)

		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), "Info Pemadaman Air Diperbarui", res.Title)
		assert.True(suite.T(), res.IsPinned)
	})

	suite.Run("Fail - Update Not Found", func() {
		suite.SetupTest()
		suite.mockNews.EXPECT().GetByID(newsID).Return(nil, errors.New("news not found"))

		svc := service.NewNewsService(suite.mockNews)
		res, err := svc.Update(newsID, input)

		assert.NotNil(suite.T(), err)
		assert.Nil(suite.T(), res)
	})
}

func (suite *NewsServiceTestSuite) TestDelete() {
	newsID := int64(8)

	suite.Run("Success - Deleted", func() {
		suite.SetupTest()
		suite.mockNews.EXPECT().GetByID(newsID).Return(&entity.RegionalNews{ID: newsID}, nil)
		suite.mockNews.EXPECT().Delete(newsID).Return(nil)

		svc := service.NewNewsService(suite.mockNews)
		err := svc.Delete(newsID)

		assert.Nil(suite.T(), err)
	})

	suite.Run("Fail - Delete Not Found", func() {
		suite.SetupTest()
		suite.mockNews.EXPECT().GetByID(newsID).Return(nil, errors.New("news not found"))

		svc := service.NewNewsService(suite.mockNews)
		err := svc.Delete(newsID)

		assert.NotNil(suite.T(), err)
	})
}

func TestNewsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NewsServiceTestSuite))
}
