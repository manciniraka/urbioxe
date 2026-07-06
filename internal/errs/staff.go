package errs

import "errors"

var (

	// Staff
	ErrStaffNotFound      = errors.New("staff not found")
	ErrStaffInactive      = errors.New("staff is inactive")
	ErrStaffAlreadyExists = errors.New("staff already exists")
)
