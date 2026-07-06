package errs

import "errors"

var (
	// Authentication
	ErrInvalidCredential = errors.New("invalid email or password")
	ErrEmailRegistered   = errors.New("email already registered")
)
