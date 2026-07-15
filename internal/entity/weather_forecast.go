package entity

import "time"

type WeatherForecast struct {
	ID uint `gorm:"primaryKey"`
	DistrictID uint `gorm:"not null"`

	ForecastTime time.Time `gorm:"not null"`
	Temperature float64 `gorm:"not null"`
	Humidity int `gorm:"not null"`
	Weather string `gorm:"size:100;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	District *District `gorm:"foreignKey:DistrictID"`
}