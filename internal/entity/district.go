package entity

import "time"

type District struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
	BMKGADM4Code string `json:"bmkg_adm4_code" gorm:"column:bmkg_adm4_code;type:varchar(20);not null"`

	Name string `json:"name" gorm:"type:varchar(100);not null"`

	IsActive bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
