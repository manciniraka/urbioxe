package service

import (
	"time"

	"github.com/manciniraka/urbioxe/internal/dto"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/repository"
)

type WaterService interface {
	CreateWaterStatus(userID uint, input dto.CreateWaterStatusInput) (*dto.WaterStatusResponse, error)
	GetAllWaterStatus() (*dto.WaterListResult, error)
	GetWaterStatusByDistrictID(districtID uint) (*dto.WaterStatusResponse, error)
	GetMyWaterStatus(userID uint) (*dto.WaterStatusResponse, error)
	GetWaterStatusHistories() ([]dto.WaterHistoryResponse, error)
}


type waterService struct {
	waterRepo repository.WaterRepository
	districtRepo repository.DistrictRepository
	userRepo repository.UserRepository
}

func NewWaterService(
	waterRepo repository.WaterRepository,
	districtRepo repository.DistrictRepository,
	userRepo repository.UserRepository,
) WaterService {
	return &waterService{
		waterRepo: waterRepo,
		districtRepo: districtRepo,
		userRepo: userRepo,
	}
}

func (ws *waterService) CreateWaterStatus(userID uint, input dto.CreateWaterStatusInput) (*dto.WaterStatusResponse, error) {
	district, err := ws.districtRepo.GetByID(
		input.DistrictID,
	)
	if err != nil {
		return nil, err
	}

	location, _ := time.LoadLocation("Asia/Jakarta")

	var startedAt time.Time
	var estimatedRecoveryAt time.Time
	var estimatedDuration int

	if input.Status != entity.WaterStatusNormal {
		startedAt, err = time.ParseInLocation(
			"2006-01-02 15:04",
			input.StartedAt,
			location,
		)
		if err != nil {
			return nil, err
		}

		now := time.Now().In(location)

		if startedAt.Before(now) {
			return nil, errs.ErrStartedAtBeforeNow
		}

		estimatedDuration = input.EstimatedDuration

		estimatedRecoveryAt = startedAt.Add(
			time.Duration(estimatedDuration) * time.Hour,
		)
	}

	waterStatus := entity.WaterStatus{
		DistrictID:           input.DistrictID,
		Status:               input.Status,
		StartedAt:            startedAt,
		EstimatedDuration:    input.EstimatedDuration,
		EstimatedRecoveryAt:  estimatedRecoveryAt,
		Reason:               input.Reason,
		CreatedBy:            userID,
	}

	err = ws.waterRepo.CreateWaterStatus(&waterStatus)
	if err != nil {
		return nil, err
	}

	response := dto.WaterStatusResponse{
		District: district.Name,
		Status: waterStatus.Status,
		StartedAt: waterStatus.
			StartedAt.
			In(location).
			Format("02 Jan 2006 15:04 WIB"),
		EstimatedRecovery: waterStatus.
			EstimatedRecoveryAt.
			In(location).
			Format("02 Jan 2006 15:04 WIB"),
		EstimatedDuration: waterStatus.EstimatedDuration,
		Reason: waterStatus.Reason,
	}

	return &response, nil
}

func (ws *waterService) GetAllWaterStatus() (*dto.WaterListResult, error) {
	statuses, err := ws.waterRepo.GetLatestAllWaterStatus()
	if err != nil {
		return nil, err
	}

	location, _ := time.LoadLocation(
		"Asia/Jakarta",
	)

	var result dto.WaterListResult

	for _, status := range statuses {

		item := dto.WaterListItem{
			District: status.District.Name,
			Status:   status.Status,
		}

		if status.Status != entity.WaterStatusNormal {
			startedAt := status.
				StartedAt.
				In(location).
				Format("02 Jan 2006 15:04 WIB")

			item.StartedAt = &startedAt
		}

		result.Data = append(
			result.Data,
			item,
		)

		if result.Metadata.UpdatedAt == "" {
			result.Metadata.UpdatedAt = status.CreatedAt.
			In(location).
			Format(
				"02 Jan 2006 15:04 WIB",
			)
		}
	}

	return &result, nil
}

func (ws *waterService) GetWaterStatusByDistrictID(districtID uint) (*dto.WaterStatusResponse, error) {
	status, err := ws.waterRepo.GetLatestWaterStatusByDistrictID(districtID)
	if err != nil {
		return nil, err
	}

	location, _ := time.LoadLocation(
		"Asia/Jakarta",
	)

	response := dto.WaterStatusResponse{
		District: status.District.Name,
		Status: status.Status,
		StartedAt: status.StartedAt.
			In(location).
			Format("02 Jan 2006 15:04 WIB"),
		EstimatedRecovery: status.
			EstimatedRecoveryAt.
			In(location).
			Format("02 Jan 2006 15:04 WIB"),
		EstimatedDuration: status.EstimatedDuration,
		Reason: status.Reason,
	}

	return &response, nil
}

func (ws *waterService) GetMyWaterStatus(userID uint) (*dto.WaterStatusResponse, error) {
	user, err := ws.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user.HomeDistrictID == nil {
		return nil, errs.ErrDistrictNotFound
	}

	return ws.GetWaterStatusByDistrictID(
		*user.HomeDistrictID,
	)
}

func (ws *waterService) GetWaterStatusHistories() ([]dto.WaterHistoryResponse,	error) {
	histories, err := ws.waterRepo.GetWaterStatusHistories()
	if err != nil {
		return nil, err
	}

	location, _ := time.LoadLocation(
		"Asia/Jakarta",
	)

	var responses []dto.WaterHistoryResponse

	for _, history := range histories {
		response := dto.WaterHistoryResponse{
			District: history.District.Name,
			Status: history.Status,
			Reason: history.Reason,
			CreatedBy: history.User.Name,
			CreatedAt: history.
				CreatedAt.
				In(location).
				Format("02 Jan 2006 15:04 WIB"),
		}

		if history.Status != entity.WaterStatusNormal {
			startedAt := history.
				StartedAt.
				In(location).
				Format("02 Jan 2006 15:04 WIB")
			estimatedRecovery := history.
				EstimatedRecoveryAt.
				In(location).
				Format("02 Jan 2006 15:04 WIB")
			estimatedDuration := history.EstimatedDuration
			response.StartedAt = &startedAt
			response.EstimatedRecovery = &estimatedRecovery
			response.EstimatedDuration = &estimatedDuration
		}

		responses = append(
			responses,
			response,
		)
	}

	return responses, nil
}