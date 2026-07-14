package errs

import "errors"

var (

	// Weather
	ErrWeatherNotFound            = errors.New("weather data not found")
	ErrWeatherSyncFailed          = errors.New("failed to synchronize weather data")
	ErrDistrictCoordinateNotFound = errors.New("district coordinate not found")
)
