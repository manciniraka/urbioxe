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

type ReportStatus string

const (
	StatusPending    ReportStatus = "pending"
	StatusVerified   ReportStatus = "verified"
	StatusAssigned   ReportStatus = "assigned"
	StatusInProgress ReportStatus = "in_progress"
	StatusResolved   ReportStatus = "resolved"
	StatusRejected   ReportStatus = "rejected"
)

type ReportPriority string

const (
	PriorityLow      ReportPriority = "low"
	PriorityMedium   ReportPriority = "medium"
	PriorityHigh     ReportPriority = "high"
	PriorityCritical ReportPriority = "critical"
)
