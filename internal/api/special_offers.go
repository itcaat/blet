package api

import (
	"fmt"

	"github.com/itcaat/blet/internal/models"
)

func (c *Client) GetSpecialOffers(origin string) (models.SpecialOffersResponse, error) {
	const url = "https://api.travelpayouts.com/aviasales/v3/get_special_offers"

	var result models.SpecialOffersResponse

	resp, err := c.resty.R().
		SetQueryParams(map[string]string{
			"origin": origin,
		}).
		SetResult(&result).
		Get(url)

	if err != nil {
		return result, err
	}

	if !result.Success {
		return result, fmt.Errorf("API error. HTTP: %s. Body: %s", resp.Status(), resp.Body())
	}

	return result, nil
}
