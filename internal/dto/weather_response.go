package dto

type CurrentWeatherResponse struct {
	Weather     string  `json:"weather"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	ValidUntil  string  `json:"valid_until"`
}

type ForecastTimelineResponse struct {
	Time        string  `json:"time"`
	Weather     string  `json:"weather"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
}

type WeatherSummaryResponse struct {
	District     string                     `json:"district"`
	ForecastDate string                     `json:"forecast_date"`
	Current      CurrentWeatherResponse     `json:"current"`
	Forecast     []ForecastTimelineResponse `json:"forecast"`
}

type WeatherListResponse struct {
	District    string  `json:"district"`
	Weather     string  `json:"weather"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
}

type SyncForecastResponse struct {
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Total   int `json:"total"`
}

type WeatherListMetadata struct {
	ForecastTime string `json:"forecast_time"`
}

type WeatherListResult struct {
	Metadata WeatherListMetadata   `json:"metadata"`
	Data     []WeatherListResponse `json:"data"`
}
