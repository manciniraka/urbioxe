package bmkg

import (
	"net/url"
	"path"
)

func (c *Client) forecastURL(adm4 string) (string, error) {
	base, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return "", err
	}

	base.Path = path.Join(base.Path, c.config.ForecastEndpoint)

	query := base.Query()
	query.Set("adm4", adm4)
	base.RawQuery = query.Encode()

	return base.String(), nil
}