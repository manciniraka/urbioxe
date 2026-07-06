package errs

import "errors"

var (

	// Emergency Contact
	ErrEmergencyContactNotFound = errors.New("emergency contact not found")
	ErrEmergencyContactInactive = errors.New("emergency contact is inactive")
)
