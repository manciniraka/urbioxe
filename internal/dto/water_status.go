package dto

type CreateWaterStatusInput struct {
	DistrictID        uint   `json:"district_id" validate:"required"`
	Status            string `json:"status" validate:"required"`
	StartedAt         string `json:"started_at"`
	EstimatedDuration int    `json:"estimated_duration"`
	Reason            string `json:"reason"`
}

type WaterStatusResponse struct {
	District          string `json:"district"`
	Status            string `json:"status"`
	StartedAt         string `json:"started_at,omitempty"`
	EstimatedRecovery string `json:"estimated_recovery,omitempty"`
	EstimatedDuration int    `json:"estimated_duration,omitempty"`
	Reason            string `json:"reason,omitempty"`
}

type WaterListItem struct {
	District  string  `json:"district"`
	Status    string  `json:"status"`
	StartedAt *string `json:"started_at,omitempty"`
}

type WaterListMetadata struct {
	UpdatedAt string `json:"updated_at"`
}

type WaterListResult struct {
	Metadata WaterListMetadata `json:"metadata"`
	Data     []WaterListItem   `json:"data"`
}

type WaterHistoryResponse struct {
	District          string  `json:"district"`
	Status            string  `json:"status"`
	StartedAt         *string `json:"started_at"`
	EstimatedRecovery *string `json:"estimated_recovery"`
	EstimatedDuration *int    `json:"estimated_duration"`
	Reason            string  `json:"reason"`
	CreatedBy         string  `json:"created_by"`
	CreatedAt         string  `json:"created_at"`
}