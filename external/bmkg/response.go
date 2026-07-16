package bmkg

type ForecastResponse struct {
	Location Location       `json:"lokasi"`
	Data     []ForecastData `json:"data"`
}

type ForecastData struct {
	Location Location    `json:"lokasi"`
	Weather  [][]Weather `json:"cuaca"`
}

type Location struct {
	ADM1 string `json:"adm1"`
	ADM2 string `json:"adm2"`
	ADM3 string `json:"adm3"`
	ADM4 string `json:"adm4"`

	Province string `json:"provinsi"`
	Regency  string `json:"kotkab"`
	District string `json:"kecamatan"`
	Village  string `json:"desa"`
}

type Weather struct {
	LocalDateTime string  `json:"local_datetime"`
	Temperature   float64 `json:"t"`
	Humidity      int     `json:"hu"`
	Weather       string  `json:"weather_desc"`
}
