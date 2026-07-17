package entity

import "time"

const (
	WaterStatusNormal      = "Normal"
	WaterStatusMaintenance = "Maintenance"
	WaterStatusLimited     = "Limited"
	WaterStatusDisrupted   = "Disrupted"
)

type WaterStatus struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
	DistrictID uint `json:"district_id"`
	Status string `json:"status" gorm:"type:varchar(30);not null"`
	StartedAt time.Time `json:"started_at"`
	EstimatedDuration int `json:"estimated_duration"`
	EstimatedRecoveryAt time.Time `json:"estimated_recovery_at"`
	Reason string `json:"reason" gorm:"type:text"`
	CreatedBy uint `json:"created_by"`
	
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:CreatedBy"`
	District District `gorm:"foreignKey:DistrictID"`
}