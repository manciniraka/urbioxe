package errs

import "errors"

var (

	// Department
	ErrDepartmentNotFound = errors.New("department not found")
	ErrDepartmentInactive = errors.New("department is inactive")
)
