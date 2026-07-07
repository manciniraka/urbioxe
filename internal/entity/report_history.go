package entity

import "time"

type ReportHistory struct {
	ID         int64        `gorm:"primaryKey;autoIncrement"`
	ReportID   int64        `gorm:"not null"`
	Status     ReportStatus `gorm:"type:report_status;not null"`
	Notes      string       `gorm:"type:text"`
	IsInternal bool         `gorm:"default:false"`
	ActorID    *int64       `gorm:"default:null"`
	CreatedAt  time.Time    `gorm:"default:NOW()"`
}
