package errs

import "errors"

var (

	// Report
	ErrReportNotFound         = errors.New("report not found")
	ErrReportForbidden        = errors.New("you are not allowed to access this report")
	ErrReportAlreadyAssigned  = errors.New("report already assigned")
	ErrReportNotAssigned      = errors.New("report has not been assigned to a staff")
	ErrReportAlreadyResolved  = errors.New("report already resolved")
	ErrReportAlreadyInProcess = errors.New("report already in process")
	ErrReportAlreadyRejected  = errors.New("report already rejected")
	ErrReportShouldInProcess  = errors.New("only in process report can be updated")
	ErrReportAssignForbidden  = errors.New("you are not allowed to assign this report")
	ErrReportUpdateForbidden  = errors.New("you are not allowed to update status this report")
)
