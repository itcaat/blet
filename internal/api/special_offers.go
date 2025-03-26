package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/itcaat/blet/internal/models"
)

const (
	getSpecialOffersURL = "https://api.travelpayouts.com/aviasales/v3/get_special_offers"
)

func GetSpecialOffers(origin, token string) (models.SpecialOffersResponse, error) {
	var client = resty.New()
	var result models.SpecialOffersResponse

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origin": origin,
			"token":  token,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get(getSpecialOffersURL)

	if err != nil {
		return result, err
	}

	if !result.Success {
		return result, fmt.Errorf("⚠️ API не вернул успешный ответ. HTTP: %s. Body: %s", resp.Status(), resp.Body())
	}

	return result, nil
}
