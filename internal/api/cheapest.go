package api

import (
	"fmt"

	"github.com/itcaat/blet/internal/models"
)

func (c *Client) GetCheapest(origin, destination, one_way string) (models.PriceForDatesResponse, error) {
	const apiUrl = "https://api.travelpayouts.com/aviasales/v3/prices_for_dates"

	var result models.PriceForDatesResponse

	resp, err := c.resty.R().
		SetQueryParams(map[string]string{
			"origin":      origin,
			"destination": destination,
			"one_way":     one_way,
			"limit":       "100",
		}).
		SetResult(&result).
		Get(apiUrl)

	if err != nil {
		return result, err
	}

	if !result.Success {
		return result, fmt.Errorf("⚠️ API не вернул успешный ответ. HTTP: %s. Body: %s. Request: %s", resp.Status(), resp.Body(), resp.Request.URL)
	}

	return result, nil
}
