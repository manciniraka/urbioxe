package service_test

import (
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/manciniraka/urbioxe/external/mailjet"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/repository"
	service "github.com/manciniraka/urbioxe/internal/service"
	"github.com/manciniraka/urbioxe/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReportServiceTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	mockRepo      *mocks.MockReportRepository
	mockStaffRepo *mocks.MockStaffRepository
	mockCldSvc    *mocks.MockCloudinaryService
	mockServer    *httptest.Server
	mailer        *mailjet.Client
}

func (suite *ReportServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())

	suite.mockRepo = mocks.NewMockReportRepository(suite.ctrl)
	suite.mockStaffRepo = mocks.NewMockStaffRepository(suite.ctrl)
	suite.mockCldSvc = mocks.NewMockCloudinaryService(suite.ctrl)

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
}

func (suite *ReportServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
	if suite.mockServer != nil {
		suite.mockServer.Close()
	}
}

func TestReportServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReportServiceTestSuite))
}

func (suite *ReportServiceTestSuite) TestGetAllReports() {
	tests := []struct {
		name       string
		param      service.GetReportsParam
		mockRepoFn func(p service.GetReportsParam)
		verifyFn   func(res *service.ReportListResponse, err error)
	}{
		{
			name: "Success - Get reports as Citizen (Uses default page and limit)",
			param: service.GetReportsParam{
				UserID: 10,
				Role:   "citizen",
				Page:   0,
				Limit:  0,
			},
			mockRepoFn: func(p service.GetReportsParam) {
				expectedFilter := repository.ReportFilter{
					UserID: 10,
					Role:   "citizen",
					Limit:  10,
					Offset: 0,
				}
				mockReports := []entity.Report{
					{ID: 1, Title: "Jalan Rusak"},
				}
				suite.mockRepo.EXPECT().
					FindAll(expectedFilter).
					Return(mockReports, int64(1), nil)
			},
			verifyFn: func(res *service.ReportListResponse, err error) {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), res)
				assert.Equal(suite.T(), 1, res.Page)
				assert.Equal(suite.T(), 10, res.Limit)
				assert.Equal(suite.T(), int64(1), res.TotalData)
				assert.Equal(suite.T(), 1, res.TotalPage)
				assert.Len(suite.T(), res.Data, 1)
			},
		},
		{
			name: "Success - Get reports as Admin/Staff (Fetches Staff ID first)",
			param: service.GetReportsParam{
				UserID: 20,
				Role:   "department_admin",
				Page:   2,
				Limit:  5,
			},
			mockRepoFn: func(p service.GetReportsParam) {
				dummyStaff := &entity.StaffProfile{
					ID:     77,
					UserID: 20,
				}
				suite.mockStaffRepo.EXPECT().
					FindStaffByUserID(uint(20)).
					Return(dummyStaff, nil)

				expectedFilter := repository.ReportFilter{
					UserID: 77,
					Role:   "department_admin",
					Limit:  5,
					Offset: 5,
				}
				mockReports := []entity.Report{
					{ID: 2, Title: "Sampah Menumpuk"},
				}
				suite.mockRepo.EXPECT().
					FindAll(expectedFilter).
					Return(mockReports, int64(6), nil)
			},
			verifyFn: func(res *service.ReportListResponse, err error) {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), res)
				assert.Equal(suite.T(), 2, res.Page)
				assert.Equal(suite.T(), 2, res.TotalPage)
				assert.Equal(suite.T(), int64(6), res.TotalData)
			},
		},
		{
			name: "Fail - Staff Profile Not Found",
			param: service.GetReportsParam{
				UserID: 99,
				Role:   "officer",
			},
			mockRepoFn: func(p service.GetReportsParam) {
				suite.mockStaffRepo.EXPECT().
					FindStaffByUserID(uint(99)).
					Return(nil, errors.New("staff not found"))
			},
			verifyFn: func(res *service.ReportListResponse, err error) {
				assert.NotNil(suite.T(), err)
				assert.Nil(suite.T(), res)
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn(tc.param)

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)
			res, err := svc.GetAllReports(tc.param)

			tc.verifyFn(res, err)
		})
	}
}

func (suite *ReportServiceTestSuite) TestGetReportByID() {
	reportID := uint(100)
	currentUserID := uint(5)

	tests := []struct {
		name          string
		role          string
		mockRepoFn    func()
		expectedError string
	}{
		{
			name: "Success - Department Admin views detail report",
			role: "department_admin",
			mockRepoFn: func() {
				mockData := &entity.Report{
					ID:    reportID,
					Title: "Banjir Luapan Drainase",
					User:  &entity.User{ID: currentUserID, Name: "Ahmad", Password: ""},
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "department_admin").Return(mockData, nil)
			},
			expectedError: "",
		},
		{
			name: "Fail - Report Not Found",
			role: "super_admin",
			mockRepoFn: func() {
				suite.mockRepo.EXPECT().FindByID(reportID, "super_admin").Return(nil, errors.New("laporan tidak ditemukan"))
			},
			expectedError: "laporan tidak ditemukan",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			result, err := svc.GetReportByID(reportID, currentUserID, tc.role)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
				assert.Nil(suite.T(), result)
			} else {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), result)
				assert.Equal(suite.T(), reportID, result.ID)
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestVerifyReport() {
	adminID := uint(1)
	reportID := uint(10)

	tests := []struct {
		name          string
		role          string
		statusBefore  entity.ReportStatus
		mockRepoFn    func()
		expectedError string
	}{
		{
			name:         "Success - Department Admin verifies pending report",
			role:         "department_admin",
			statusBefore: entity.StatusPending,
			mockRepoFn: func() {
				reportData := &entity.Report{
					ID:     reportID,
					Title:  "Jalan Berlubang",
					Status: entity.StatusPending,
					User:   &entity.User{Email: "warga@mail.com", Name: "Rian"},
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "department_admin").Return(reportData, nil)
				suite.mockRepo.EXPECT().UpdateStatusWithHistory(reportID, entity.StatusVerified, adminID, gomock.Any(), false).Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Fail - Unauthorized Role (Officer)",
			role:          "officer",
			statusBefore:  entity.StatusPending,
			mockRepoFn:    func() {},
			expectedError: "you are not allowed to update status this report",
		},
		{
			name:         "Fail - Report is already In Progress",
			role:         "super_admin",
			statusBefore: entity.StatusInProgress,
			mockRepoFn: func() {
				reportData := &entity.Report{ID: reportID, Status: entity.StatusInProgress}
				suite.mockRepo.EXPECT().FindByID(reportID, "super_admin").Return(reportData, nil)
			},
			expectedError: "only in pending report can be updated",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			input := service.UpdateStatusReportInput{Notes: "Laporan terverifikasi valid"}
			err := svc.VerifyReport(reportID, adminID, tc.role, input)

			time.Sleep(15 * time.Millisecond)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
			} else {
				assert.Nil(suite.T(), err)
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestUpdatePriority() {
	adminID := uint(2)
	reportID := uint(42)

	tests := []struct {
		name          string
		role          string
		inputPriority entity.ReportPriority
		mockRepoFn    func()
		expectedError string
	}{
		{
			name:          "Success - Super Admin updates priority to Critical",
			role:          "super_admin",
			inputPriority: "critical",
			mockRepoFn: func() {
				suite.mockRepo.EXPECT().
					FindByID(reportID, "super_admin").
					Return(&entity.Report{ID: reportID, Status: entity.StatusPending}, nil)

				suite.mockRepo.EXPECT().
					UpdatePriorityWithHistory(
						reportID,
						entity.StatusPending,
						entity.ReportPriority("critical"),
						adminID,
						gomock.Any(),
					).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Fail - Officer tries to change priority",
			role:          "officer",
			inputPriority: "high",
			mockRepoFn:    func() {},
			expectedError: "forbidden",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			input := service.UpdatePriorityInput{Priority: tc.inputPriority, Notes: "Segera ditindak"}
			err := svc.UpdatePriority(reportID, adminID, tc.role, input)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
			} else {
				assert.Nil(suite.T(), err)
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestStartReport() {
	officerUserID := uint(3)
	staffID := uint(30)
	reportID := uint(55)

	tests := []struct {
		name          string
		role          string
		notes         string
		mockRepoFn    func()
		expectedError string
	}{
		{
			name:  "Success - Assigned Officer starts working on the report",
			role:  "officer",
			notes: "Mulai pengerjaan jalan berlubang",
			mockRepoFn: func() {
				currentStaffID := staffID
				reportData := &entity.Report{
					ID:              reportID,
					Status:          entity.StatusAssigned,
					AssignedStaffID: &currentStaffID,
					User:            &entity.User{Email: "warga@mail.com", Name: "Rian"},
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "officer").Return(reportData, nil)

				suite.mockStaffRepo.EXPECT().
					FindStaffByUserID(officerUserID).
					Return(&entity.StaffProfile{ID: staffID}, nil)

				suite.mockRepo.EXPECT().
					UpdateStatusWithHistory(reportID, entity.StatusInProgress, officerUserID, "Mulai pengerjaan jalan berlubang", false).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:  "Fail - Forbidden if another officer tries to start it",
			role:  "officer",
			notes: "Coba handle laporan orang lain",
			mockRepoFn: func() {
				wrongStaffID := uint(99)
				reportData := &entity.Report{
					ID:              reportID,
					Status:          entity.StatusAssigned,
					AssignedStaffID: &wrongStaffID,
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "officer").Return(reportData, nil)

				suite.mockStaffRepo.EXPECT().
					FindStaffByUserID(officerUserID).
					Return(&entity.StaffProfile{ID: staffID}, nil)
			},
			expectedError: "forbidden",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			err := svc.StartReport(reportID, officerUserID, tc.role, tc.notes)

			time.Sleep(15 * time.Millisecond)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
			} else {
				assert.Nil(suite.T(), err)
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestResolveReport() {
	officerUserID := uint(3)
	staffID := uint(30)
	reportID := uint(55)

	mockFileHeader := &multipart.FileHeader{
		Filename: "resolved_evidence.jpg",
		Size:     2048,
	}
	validFiles := []*multipart.FileHeader{mockFileHeader}

	tests := []struct {
		name          string
		role          string
		notes         string
		files         []*multipart.FileHeader
		mockRepoFn    func()
		expectedError string
	}{
		{
			name:  "Success - Officer resolves in-progress report with evidence",
			role:  "officer",
			notes: "Pekerjaan perbaikan jalan selesai",
			files: validFiles,
			mockRepoFn: func() {
				currentStaffID := staffID
				reportData := &entity.Report{
					ID:              reportID,
					Status:          entity.StatusInProgress,
					AssignedStaffID: &currentStaffID,
					User:            &entity.User{Email: "warga@mail.com", Name: "Rian"},
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "officer").Return(reportData, nil)

				suite.mockStaffRepo.EXPECT().
					FindStaffByUserID(officerUserID).
					Return(&entity.StaffProfile{ID: staffID}, nil)

				suite.mockCldSvc.EXPECT().
					UploadImage(mockFileHeader).
					Return("https://cloudinary.com/resolved.jpg", nil)

				suite.mockRepo.EXPECT().
					ResolveReport(reportID, officerUserID, "Pekerjaan perbaikan jalan selesai", gomock.Any()).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Fail - No image uploaded",
			role:          "officer",
			notes:         "Selesai",
			files:         []*multipart.FileHeader{},
			mockRepoFn:    func() {},
			expectedError: "upload minimal 1 image",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			err := svc.ResolveReport(reportID, officerUserID, tc.role, tc.notes, tc.files)

			time.Sleep(15 * time.Millisecond)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
			} else {
				assert.Nil(suite.T(), err)
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestRejectReport() {
	adminID := uint(1)
	reportID := uint(77)

	tests := []struct {
		name          string
		role          string
		mockRepoFn    func()
		expectedError string
	}{
		{
			name: "Success - Department Admin rejects a report",
			role: "department_admin",
			mockRepoFn: func() {
				reportData := &entity.Report{
					ID:     reportID,
					Status: entity.StatusPending,
					User:   &entity.User{Email: "warga@mail.com", Name: "Rian"},
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "department_admin").Return(reportData, nil)
				suite.mockRepo.EXPECT().UpdateStatusWithHistory(reportID, entity.StatusRejected, adminID, "Foto tidak jelas", false).Return(nil)
			},
			expectedError: "",
		},
		{
			name:          "Fail - Officer cannot reject a report",
			role:          "officer",
			mockRepoFn:    func() {},
			expectedError: "you are not allowed to update status this report",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			input := service.UpdateStatusReportInput{
				Notes: "Foto tidak jelas",
			}
			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)
			err := svc.RejectReport(reportID, adminID, tc.role, input)

			time.Sleep(15 * time.Millisecond)
			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
			} else {
				assert.Nil(suite.T(), err)
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestCreateReport() {
	userID := uint(10)

	mockFileHeader := &multipart.FileHeader{
		Filename: "evidence.jpg",
		Size:     1024,
	}
	validFiles := []*multipart.FileHeader{mockFileHeader}

	tests := []struct {
		name          string
		inputTitle    string
		districtID    uint
		categoryID    uint
		files         []*multipart.FileHeader
		mockRepoFn    func()
		expectedError string
	}{
		{
			name:       "Success - Citizen creates a valid report with an image",
			inputTitle: "Jalan Rusak Parah",
			districtID: 5,
			categoryID: 1,
			files:      validFiles,
			mockRepoFn: func() {
				suite.mockCldSvc.EXPECT().
					UploadImage(mockFileHeader).
					Return("https://cloudinary.com/evidence.jpg", nil)

				suite.mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(report *entity.Report) error {
					report.ID = uint(100)
					assert.NotEmpty(suite.T(), *report.ReportNumber)
					assert.Equal(suite.T(), entity.StatusPending, report.Status)
					assert.Len(suite.T(), report.Histories, 1)
					return nil
				})
			},
			expectedError: "",
		},
		{
			name:       "Fail - Cloudinary Upload Error",
			inputTitle: "Sampah Menumpuk",
			districtID: 3,
			categoryID: 2,
			files:      validFiles,
			mockRepoFn: func() {
				suite.mockCldSvc.EXPECT().
					UploadImage(mockFileHeader).
					Return("", errors.New("cloudinary storage full"))
			},
			expectedError: "cloudinary storage full",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			input := entity.Report{
				Title:              tc.inputTitle,
				IncidentDistrictID: tc.districtID,
				CategoryID:         tc.categoryID,
				UserID:             userID,
				User: &entity.User{
					ID:    userID,
					Name:  "Rian Warga",
					Email: "rian@mail.com",
				},
				Category: &entity.Category{
					ID:   tc.categoryID,
					Name: "Infrastruktur",
				},
				AddressLandmark: "Dekat halte bus",
			}

			result, err := svc.CreateReport(input, tc.files)

			time.Sleep(15 * time.Millisecond)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
				assert.Nil(suite.T(), result)
			} else {
				assert.Nil(suite.T(), err)
				assert.NotNil(suite.T(), result)
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestAssignReport() {
	adminUserID := uint(1)
	reportID := uint(10)
	targetStaffID := uint(50)

	tests := []struct {
		name          string
		role          string
		input         service.AssignReportInput
		mockRepoFn    func()
		expectedError string
	}{
		{
			name: "Success - Admin assigns report to a staff",
			role: "department_admin",
			input: service.AssignReportInput{
				StaffID: targetStaffID,
				Notes:   "Ditugaskan ke regu lapangan",
			},
			mockRepoFn: func() {
				reportData := &entity.Report{
					ID:              reportID,
					Status:          entity.StatusVerified,
					AssignedStaffID: nil,
					User:            &entity.User{Email: "warga@mail.com", Name: "Rian"},
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "department_admin").Return(reportData, nil)

				suite.mockRepo.EXPECT().
					AssignStaff(reportID, targetStaffID, adminUserID, "Ditugaskan ke regu lapangan").
					Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Fail - StaffID is missing",
			role: "department_admin",
			input: service.AssignReportInput{
				StaffID: 0,
				Notes:   "Test",
			},
			mockRepoFn:    func() {},
			expectedError: "staff_id required",
		},
		{
			name: "Fail - Officer tries to assign report",
			role: "officer",
			input: service.AssignReportInput{
				StaffID: targetStaffID,
			},
			mockRepoFn:    func() {},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			err := svc.AssignReport(reportID, adminUserID, tc.role, tc.input)

			time.Sleep(15 * time.Millisecond)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
			} else {
				if tc.role == "officer" {
					assert.NotNil(suite.T(), err)
				} else {
					assert.Nil(suite.T(), err)
				}
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestUpdateReport() {
	userID := uint(10)
	reportID := uint(20)

	tests := []struct {
		name          string
		requesterID   uint
		status        entity.ReportStatus
		input         service.UpdateReportInput
		mockRepoFn    func(status entity.ReportStatus, reqID uint)
		expectedError string
	}{
		{
			name:        "Success - Citizen updates pending report details",
			requesterID: userID,
			status:      entity.StatusPending,
			input: service.UpdateReportInput{
				Title:           "Judul Baru Hasil Revisi Warga",
				AddressLandmark: "Patokan depan minimarket",
			},
			mockRepoFn: func(status entity.ReportStatus, reqID uint) {
				oldReport := &entity.Report{
					ID:     reportID,
					UserID: userID,
					Status: status,
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "citizen").Return(oldReport, nil)

				suite.mockRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectedError: "",
		},
		{
			name:        "Fail - Forbidden if another citizen tries to update",
			requesterID: uint(99),
			status:      entity.StatusPending,
			input: service.UpdateReportInput{
				Title: "Mencoba edit laporan orang lain",
			},
			mockRepoFn: func(status entity.ReportStatus, reqID uint) {
				oldReport := &entity.Report{
					ID:     reportID,
					UserID: userID,
					Status: status,
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "citizen").Return(oldReport, nil)
			},
			expectedError: "",
		},
		{
			name:        "Fail - Report already in process",
			requesterID: userID,
			status:      entity.StatusInProgress,
			input: service.UpdateReportInput{
				Title: "Edit laporan yang sedang dikerjakan",
			},
			mockRepoFn: func(status entity.ReportStatus, reqID uint) {
				oldReport := &entity.Report{
					ID:     reportID,
					UserID: userID,
					Status: status,
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "citizen").Return(oldReport, nil)
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn(tc.status, tc.requesterID)

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			result, err := svc.UpdateReport(reportID, tc.requesterID, tc.input)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
				assert.Nil(suite.T(), result)
			} else {
				if tc.requesterID != userID || tc.status != entity.StatusPending {
					assert.NotNil(suite.T(), err)
					assert.Nil(suite.T(), result)
				} else {
					assert.Nil(suite.T(), err)
					assert.NotNil(suite.T(), result)
				}
			}
		})
	}
}

func (suite *ReportServiceTestSuite) TestReassignReport() {
	adminUserID := uint(1)
	reportID := uint(12)
	oldStaffID := uint(50)
	newStaffID := uint(60)

	tests := []struct {
		name          string
		role          string
		input         service.ReassignReportInput
		mockRepoFn    func()
		expectedError string
	}{
		{
			name: "Success - Admin reassigns report to another staff",
			role: "super_admin",
			input: service.ReassignReportInput{
				NewStaffID: newStaffID,
				Notes:      "Petugas sebelumnya berhalangan sakit",
			},
			mockRepoFn: func() {
				prevStaff := oldStaffID
				reportData := &entity.Report{
					ID:              reportID,
					Status:          entity.StatusAssigned,
					AssignedStaffID: &prevStaff,
					User:            &entity.User{Email: "warga@mail.com", Name: "Rian"},
				}

				suite.mockRepo.EXPECT().FindByID(reportID, "super_admin").Return(reportData, nil)

				suite.mockRepo.EXPECT().
					AssignStaff(reportID, newStaffID, adminUserID, gomock.Any()).
					Return(nil)

				suite.mockRepo.EXPECT().FindByID(reportID, "super_admin").Return(reportData, nil)
			},
			expectedError: "",
		},
		{
			name: "Fail - NewStaffID is empty",
			role: "department_admin",
			input: service.ReassignReportInput{
				NewStaffID: 0,
				Notes:      "Test",
			},
			mockRepoFn:    func() {},
			expectedError: "staff id required",
		},
		{
			name: "Fail - Report not assigned yet",
			role: "super_admin",
			input: service.ReassignReportInput{
				NewStaffID: newStaffID,
			},
			mockRepoFn: func() {
				reportData := &entity.Report{
					ID:              reportID,
					Status:          entity.StatusVerified,
					AssignedStaffID: nil,
				}
				suite.mockRepo.EXPECT().FindByID(reportID, "super_admin").Return(reportData, nil)
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockRepoFn()

			svc := service.NewReportService(suite.mockRepo, suite.mockStaffRepo, suite.mockCldSvc, suite.mailer)

			err := svc.ReassignReport(reportID, adminUserID, tc.role, tc.input)

			time.Sleep(15 * time.Millisecond)

			if tc.expectedError != "" {
				assert.NotNil(suite.T(), err)
				assert.Contains(suite.T(), err.Error(), tc.expectedError)
			} else {
				if tc.name == "Fail - Report not assigned yet" {
					assert.NotNil(suite.T(), err)
				} else {
					assert.Nil(suite.T(), err)
				}
			}
		})
	}
}
