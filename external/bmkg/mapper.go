package bmkg

import "time"

const (
	localDateTimeLayout = "2006-01-02 15:04:05"
)

func MapForecastResponse(response ForecastResponse) ([]Forecast, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}

	var forecasts []Forecast

	for _, data := range response.Data {
		for _, weatherGroup := range data.Weather {
			for _, weather := range weatherGroup {
				forecastTime, err := time.ParseInLocation(
					localDateTimeLayout,
					weather.LocalDateTime,
					location,
				)
				if err != nil {
					return nil, err
				}

				forecasts = append(
					forecasts,
					Forecast{
						ForecastTime: forecastTime,
						Temperature: weather.Temperature,
						Humidity: weather.Humidity,
						Weather: weather.Weather,
					},
				)
			}
		}
	}

	return forecasts, nil
}