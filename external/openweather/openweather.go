package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/manciniraka/urbioxe/internal/errs"
)

type Config struct {
	BaseURL string
	APIKey  string
}

type Client struct {
	config     Config
	httpClient *http.Client
}

func New(cfg Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type CurrentWeather struct {
	Temperature float64
	Humidity    int
	Weather     string
	AirQuality  int
}

type weatherAPIResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`

	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
}

type airPollutionAPIResponse struct {
	List []struct {
		Main struct {
			AQI int `json:"aqi"`
		} `json:"main"`
	} `json:"list"`
}

// Temporary Hardcode district's Coordinate
var districtCoordinates = map[string]Coordinate{

	"Subang": {
		Latitude:  -6.5719,
		Longitude: 107.7520,
	},

	"Pamanukan": {
		Latitude:  -6.2846,
		Longitude: 107.8107,
	},

	"Ciasem": {
		Latitude:  -6.2355,
		Longitude: 107.8222,
	},

	"Kalijati": {
		Latitude:  -6.5396,
		Longitude: 107.4457,
	},

	"Pagaden": {
		Latitude:  -6.4773,
		Longitude: 107.7734,
	},
}

func (c *Client) getCoordinate(districtName string) (Coordinate, error) {
	coordinate, ok := districtCoordinates[districtName]
	if !ok {
		return Coordinate{}, errs.ErrDistrictCoordinateNotFound
	}

	return coordinate, nil
}

func (c *Client) GetCurrentWeather(districtName string) (*CurrentWeather, error) {
	coordinate, err := c.getCoordinate(
		districtName,
	)
	if err != nil {
		return nil, err
	}

	weather, err := c.getWeather(
		coordinate,
	)
	if err != nil {
		return nil, err
	}

	aqi, err := c.getAirQuality(
		coordinate,
	)
	if err != nil {
		return nil, err
	}

	weather.AirQuality = aqi

	return weather, nil
}

func (c *Client) getWeather(coordinate Coordinate) (*CurrentWeather, error) {

	url := fmt.Sprintf(
		"%s/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s",
		c.config.BaseURL,
		coordinate.Latitude,
		coordinate.Longitude,
		c.config.APIKey,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {

		return nil, fmt.Errorf(
			"openweather returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var response weatherAPIResponse

	if err := json.Unmarshal(
		responseBody,
		&response,
	); err != nil {
		return nil, err
	}

	weather := &CurrentWeather{
		Temperature: response.Main.Temp,
		Humidity:    response.Main.Humidity,
	}

	if len(response.Weather) > 0 {
		weather.Weather = response.Weather[0].Main
	}

	return weather, nil
}

func (c *Client) getAirQuality(coordinate Coordinate) (int, error) {
	url := fmt.Sprintf(
		"%s/data/2.5/air_pollution?lat=%f&lon=%f&appid=%s",
		c.config.BaseURL,
		coordinate.Latitude,
		coordinate.Longitude,
		c.config.APIKey,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {

		return 0, fmt.Errorf(
			"openweather returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var response airPollutionAPIResponse

	if err := json.Unmarshal(
		responseBody,
		&response,
	); err != nil {
		return 0, err
	}

	if len(response.List) == 0 {
		return 0, errs.ErrWeatherNotFound
	}

	return response.List[0].Main.AQI, nil
}
