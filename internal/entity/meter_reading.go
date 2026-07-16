package entity

import "time"

const (
	MeterReadingPending  = "Pending"
	MeterReadingApproved = "Approved"
	MeterReadingRejected = "Rejected"
)

type MeterReading struct {
	ID uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	CustomerNumber string `gorm:"not null"`
	CurrentReading int `gorm:"not null"`
	PhotoURL string `gorm:"not null"`
	Status string `gorm:"default:'Pending'"`
	
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}