package entity

import (
	"time"
)

type User struct {
	ID             uint  `json:"id" gorm:"primaryKey;autoIncrement"`
	HomeDistrictID *uint `json:"home_district_id,omitempty"`

	NIK         string   `json:"nik" gorm:"type:char(16);unique;not null"`
	Name        string   `json:"name" gorm:"not null"`
	Email       string   `json:"email" gorm:"unique;not null"`
	Password    string   `json:"password,omitempty" gorm:"not null"`
	PhoneNumber string   `json:"phone_number" gorm:"type:text"`
	Role        UserRole `json:"role" gorm:"type:user_role;default:'citizen'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
