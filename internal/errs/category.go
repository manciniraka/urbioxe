package errs

import "errors"

var (

	// Category
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryInactive = errors.New("category is inactive")
)
