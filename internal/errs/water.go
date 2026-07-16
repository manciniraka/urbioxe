package errs

import "errors"

var (

	// Weather
	ErrWaterStatusNotFound = errors.New("water status not found")
	ErrStartedAtBeforeNow   = errors.New("started_at cannot be before current time")
	ErrWaterTariffNotFound  = errors.New("water tariff not found")
)