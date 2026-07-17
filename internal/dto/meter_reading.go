package dto

type CreateMeterReadingInput struct {
	CustomerNumber string `form:"customer_number" validate:"required"`
	CurrentReading int    `form:"current_reading" validate:"required,min=0"`
}

type MeterReadingResponse struct {
	CustomerNumber string `json:"customer_number"`
	CurrentReading int    `json:"current_reading"`
	PhotoURL       string `json:"photo_url"`
	Status         string `json:"status"`
	SubmittedAt    string `json:"submitted_at"`
}

type MeterReadingHistoryResponse struct {
	ID             uint   `json:"id"`
	CustomerNumber string `json:"customer_number"`
	CurrentReading int    `json:"current_reading"`
	PhotoURL       string `json:"photo_url"`
	Status         string `json:"status"`
	SubmittedAt    string `json:"submitted_at"`
}

type GetMeterReadingsQuery struct {
	Month *int `query:"month"`
	Year  *int `query:"year"`
}