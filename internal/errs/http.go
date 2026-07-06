package errs

import (
	"errors"
	"net/http"
)

func StatusCode(err error) int {

	switch {

	// 400 Bad Request
	case errors.Is(err, ErrBadRequest),
		errors.Is(err, ErrPasswordMismatch),
		errors.Is(err, ErrSamePassword),
		errors.Is(err, ErrReportNotAssigned):
		return http.StatusBadRequest

	// 401 Unauthorized
	case errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrInvalidCredential):
		return http.StatusUnauthorized

	// 403 Forbidden
	case errors.Is(err, ErrForbidden),
		errors.Is(err, ErrReportForbidden):
		return http.StatusForbidden

	// 404 Not Found
	case errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrStaffNotFound),
		errors.Is(err, ErrDistrictNotFound),
		errors.Is(err, ErrDepartmentNotFound),
		errors.Is(err, ErrCategoryNotFound),
		errors.Is(err, ErrReportNotFound),
		errors.Is(err, ErrNewsNotFound),
		errors.Is(err, ErrEmergencyContactNotFound),
		errors.Is(err, ErrWeatherNotFound):
		return http.StatusNotFound

	// 409 Conflict
	case errors.Is(err, ErrEmailRegistered),
		errors.Is(err, ErrNIKRegistered),
		errors.Is(err, ErrStaffAlreadyExists),
		errors.Is(err, ErrReportAlreadyAssigned),
		errors.Is(err, ErrReportAlreadyResolved):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
