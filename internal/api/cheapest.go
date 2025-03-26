package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/itcaat/blet/internal/models"
)

const (
	pricesForDatesURL = "https://api.travelpayouts.com/aviasales/v3/prices_for_dates"
)

func GetCheapest(origin, destination, token string) (models.PriceForDatesResponse, error) {
	var client = resty.New()
	var result models.PriceForDatesResponse

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origin":      origin,
			"destination": destination,
			"one_way":     "true",
			"limit":       "1000",
			"token":       token,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get(pricesForDatesURL)

	if err != nil {
		return result, err
	}

	if !result.Success {
		return result, fmt.Errorf("⚠️ API не вернул успешный ответ. HTTP: %s. Body: %s", resp.Status(), resp.Body())
	}

	return result, nil
}
