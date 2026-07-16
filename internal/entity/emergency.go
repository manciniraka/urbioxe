package entity

type EmergencyContact struct {
	ID           uint
	DepartmentID *uint
	DistrictID   *uint
	Name         string
	PhoneNumber  string
	Description  *string
	IconURL      *string
	IsActive     bool

	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	District   *District   `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
}
