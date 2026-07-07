package entity

import "time"

type AttachmentType string

const (
	TypeEvidence   AttachmentType = "evidence"
	TypeResolution AttachmentType = "resolution"
)

type ReportAttachment struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	ReportID  int64          `gorm:"not null"`
	FileURL   string         `gorm:"type:text;not null"`
	Type      AttachmentType `gorm:"type:attachment_type;not null"`
	CreatedAt time.Time      `gorm:"default:NOW()"`
}
