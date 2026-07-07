package entity

type UserRole string

const (
	RoleCitizen         UserRole = "citizen"
	RoleOfficer         UserRole = "officer"
	RoleDepartmentAdmin UserRole = "department_admin"
	RoleSuperAdmin      UserRole = "super_admin"
)
