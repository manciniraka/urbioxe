package entity

type UserRole string

const (
	RoleCitizen         UserRole = "citizen"
	RoleOfficer         UserRole = "officer"
	RoleDepartmentAdmin UserRole = "department_admin"
	RoleSuperAdmin      UserRole = "super_admin"
)

type StaffPosition string

const (
	PositionFieldOfficer    StaffPosition = "field_officer"
	PositionDepartmentAdmin StaffPosition = "department_admin"
	PositionSupervisor      StaffPosition = "supervisor"
)