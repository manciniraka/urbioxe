package entity

type EmergencyContact struct {
	ID          uint
	Name        string
	PhoneNumber string
	Description *string
	IconURL     *string
	IsActive    bool
}
