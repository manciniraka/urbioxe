package validator

import (
	goValidator "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *goValidator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: goValidator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}