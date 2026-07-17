package entity

import (
	"time"
)

type Report struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement"`
	UserID             uint           `gorm:"not null"`
	CategoryID         uint           `gorm:"not null"`
	IncidentDistrictID uint           `gorm:"not null"`
	AssignedStaffID    *uint          `gorm:"default:null"`
	Title              string         `gorm:"type:varchar(255);not null"`
	Description        string         `gorm:"type:text;not null"`
	ReportNumber       *string        `gorm:"type:varchar(50);unique;default:null"`
	Latitude           *float64       `gorm:"type:double precision;default:null"`
	Longitude          *float64       `gorm:"type:double precision;default:null"`
	AddressLandmark    string         `gorm:"type:text"`
	Status             ReportStatus   `gorm:"type:report_status;default:'pending'"`
	Priority           ReportPriority `gorm:"type:report_priority;default:'medium'"`
	CreatedAt          time.Time      `gorm:"default:NOW()"`
	UpdatedAt          time.Time      `gorm:"default:NOW()"`

	User        *User              `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT"`
	Category    *Category          `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
	District    *District          `gorm:"foreignKey:IncidentDistrictID;constraint:OnDelete:RESTRICT"`
	Attachments []ReportAttachment `gorm:"foreignKey:ReportID;constraint:OnDelete:CASCADE"`
	Histories   []ReportHistory    `gorm:"foreignKey:ReportID;constraint:OnDelete:CASCADE"`
}
