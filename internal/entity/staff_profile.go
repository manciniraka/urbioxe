package entity

import "time"

type StaffProfile struct {
	ID           uint `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint `json:"user_id"`
	DepartmentID uint `json:"department_id"`

	EmployeeNumber string        `json:"employee_number" gorm:"unique;not null"`
	Position       StaffPosition `json:"position" gorm:"type:staff_position;not null"`
	JoinDate       string        `json:"join_date" validate:"required"`
	IsActive       bool          `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User       *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Department *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}
