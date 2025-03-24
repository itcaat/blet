package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/itcaat/blet/internal/models"
)

const (
	baseAviasalesURL  = "https://www.aviasales.ru"
	pricesForDatesURL = "https://api.travelpayouts.com/aviasales/v3/prices_for_dates"
	weekMatrixURL     = "https://api.travelpayouts.com/v2/prices/week-matrix"
)

var client = resty.New()

func GetCheapest(origin, token string) (models.PriceForDatesResponse, error) {
	var result models.PriceForDatesResponse

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origin":  origin,
			"one_way": "true",
			"limit":   "500",
			"token":   token,
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

func GetWeekPrices(origin, destination, depart, back, token string) (*models.WeekMatrixResponse, error) {
	var result models.WeekMatrixResponse

	fmt.Printf("Запрашиваю данные...: %s → %s %s - %s\n", origin, destination, depart, back)

	params := map[string]string{
		"origin":      origin,
		"destination": destination,
		"token":       token,
		"depart_date": depart,
	}

	if back != "" {
		params["return_date"] = back
	}

	resp, err := client.R().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get(weekMatrixURL)

	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("API error: %s. Body: %s. Url: %s", resp.Status(), resp.Body(), resp.Request.URL)
	}

	return &result, nil
}
