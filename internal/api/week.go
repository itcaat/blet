package api

import (
	"fmt"

	"github.com/itcaat/blet/internal/models"
)

func (c *Client) GetWeekPrices(origin, destination, depart, back string) (models.WeekMatrixResponse, error) {
	const apiUrl = "https://api.travelpayouts.com/v2/prices/week-matrix"
	var result models.WeekMatrixResponse

	fmt.Printf("Запрашиваю данные...: %s → %s %s - %s\n", origin, destination, depart, back)

	params := map[string]string{
		"origin":      origin,
		"destination": destination,
		"depart_date": depart,
	}

	if back != "" {
		params["return_date"] = back
	}

	resp, err := c.resty.R().
		SetQueryParams(params).
		SetResult(&result).
		Get(apiUrl)

	if err != nil {
		return result, err
	}

	if !result.Success {
		return result, fmt.Errorf("API error: %s. Body: %s. Url: %s", resp.Status(), resp.Body(), resp.Request.URL)
	}

	return result, nil
}
