package errs

import "errors"

var (

	// District
	ErrDistrictNotFound = errors.New("district not found")
	ErrDistrictInactive = errors.New("district is inactive")
)
