package errs

import "errors"

var (

	// User
	ErrUserNotFound      = errors.New("user not found")
	ErrNIKRegistered     = errors.New("nik already registered")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrPasswordMismatch  = errors.New("old password is incorrect")
	ErrSamePassword      = errors.New("new password must be different from old password")
	ErrHomeDistrictNotSet = errors.New("home district has not been set")
)
