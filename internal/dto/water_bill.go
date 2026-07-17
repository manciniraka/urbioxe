package dto

type BillSimulationRequest struct {
	TariffGroup string `json:"tariff_group" validate:"required"`
	Usage       int    `json:"usage" validate:"required,min=0"`
}

type BillBreakdown struct {
	Range      string `json:"range"`
	Usage      int    `json:"usage"`
	PricePerM3 int    `json:"price_per_m3"`
	Subtotal   int    `json:"subtotal"`
}

type BillSimulationResponse struct {
	TariffGroup   string          `json:"tariff_group"`
	TariffName    string          `json:"tariff_name"`
	Usage         int             `json:"usage"`
	MinimumUsage  int             `json:"minimum_usage"`
	BillableUsage int             `json:"billable_usage"`
	TotalBill     int             `json:"total_bill"`
	Breakdown     []BillBreakdown `json:"breakdown"`
}