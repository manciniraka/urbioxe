package entity

import "time"

type ReportHistory struct {
	ID         uint         `gorm:"primaryKey;autoIncrement"`
	ReportID   uint         `gorm:"not null"`
	Status     ReportStatus `gorm:"type:report_status;not null"`
	Notes      string       `gorm:"type:text"`
	IsInternal bool         `gorm:"default:false"`
	ActorID    *uint        `gorm:"default:null"`
	CreatedAt  time.Time    `gorm:"default:NOW()"`
}
