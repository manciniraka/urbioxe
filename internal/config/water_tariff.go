package config

import "github.com/manciniraka/urbioxe/internal/errs"

type TariffBlock struct {
	From int
	To   int

	PricePerM3 int
}

type WaterTariff struct {
	Code string
	Name string

	MinimumUsage int

	Blocks []TariffBlock
}

var WaterTariffs = []WaterTariff{
	{
		Code: "RT3",
		Name: "Rumah Tangga III",
		MinimumUsage: 10,
		Blocks: []TariffBlock{
			{
				From: 1,
				To:   10,
				PricePerM3: 6000,
			},
			{
				From: 11,
				To:   20,
				PricePerM3: 8000,
			},
			{
				From: 21,
				To:   30,
				PricePerM3: 10000,
			},
			{
				From: 31,
				To:   -1,
				PricePerM3: 12000,
			},
		},
	},
}

func GetWaterTariff(code string) (*WaterTariff, error) {
	for _, tariff := range WaterTariffs {
		if tariff.Code == code {
			return &tariff, nil
		}
	}

	return nil, errs.ErrWaterTariffNotFound
}

