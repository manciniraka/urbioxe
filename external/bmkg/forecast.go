package bmkg

import "time"

type Forecast struct {
	ForecastTime time.Time
	Temperature  float64
	Humidity     int
	Weather      string
}
