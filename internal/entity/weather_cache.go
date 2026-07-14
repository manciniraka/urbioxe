package entity

import "time"

type WeatherCache struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	DistrictID  uint    `json:"district_id"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	Weather     string  `json:"weather"`
	AirQuality  int     `json:"air_quality"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	
	District District `json:"district,omitempty" gorm:"foreignKey:DistrictID"`
}