package entity

import "time"

type AttachmentType string

const (
	TypeEvidence   AttachmentType = "evidence"
	TypeResolution AttachmentType = "resolution"
)

type ReportAttachment struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	ReportID  uint           `gorm:"not null"`
	FileURL   string         `gorm:"type:text;not null"`
	Type      AttachmentType `gorm:"type:attachment_type;not null"`
	CreatedAt time.Time      `gorm:"default:NOW()"`
}
