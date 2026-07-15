package bmkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/manciniraka/urbioxe/internal/logger"
)

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

func (c *Client) GetForecast(adm4 string) ([]Forecast, error) {
	url, err := c.forecastURL(adm4)
	if err != nil {
		return nil, err
	}
	logger.Log.Info(
		"calling BMKG API",
		"url", url,
		"adm4", adm4,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"bmkg returned status %d",
			resp.StatusCode,
		)
	}

	var forecastResponse ForecastResponse

	if err := json.NewDecoder(resp.Body).Decode(
		&forecastResponse,
	); err != nil {
		return nil, err
	}

	return MapForecastResponse(
		forecastResponse,
	)
}