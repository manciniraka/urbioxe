package errs

import "errors"

var (

	// Report
	ErrReportNotFound        = errors.New("report not found")
	ErrReportForbidden       = errors.New("you are not allowed to access this report")
	ErrReportAlreadyAssigned = errors.New("report already assigned")
	ErrReportNotAssigned     = errors.New("report has not been assigned to a staff")
	ErrReportAlreadyResolved = errors.New("report already resolved")
)
