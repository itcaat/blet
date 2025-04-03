package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/itcaat/blet/internal/models"
)

const (
	partnerLinkShorterURL = "https://api.travelpayouts.com/links/v1/create"
)

func GetShortUrl(url, token string) (models.ShortLinksResponse, error) {
	var client = resty.New()

	// {
	// 	"trs": 400658,
	// 	"marker": 616825,
	// 	"shorten": true,
	// 	"links": [
	// 		{
	// 			"url": "https://www.aviasales.ru/search/LED1201KUF1?t=S717682492001768330500001295LEDDMEKUF_6f72f9df6f61f48624b3183cbc36d313_7992&search_date=26032025&expected_price_uuid=dae5e307-595b-4841-9bfa-88ee28e5ce01&expected_price_source=share&expected_price_currency=rub&expected_price=7966"
	// 		}
	// 	]
	//  }var client = resty.New()

	var result models.ShortLinksResponse

	resp, err := client.R().
		SetBody(map[string]interface{}{
			"trs":     400658,
			"marker":  616825,
			"shorten": true,
			"links": []map[string]interface{}{
				{
					"url": url,
				},
			},
		}).
		SetHeader("Accept", "application/json").
		SetHeader("X-Access-Token", token).
		SetResult(&result).
		Post(partnerLinkShorterURL)

	if err != nil {
		return result, err
	}

	if result.Status != "success" {
		return result, fmt.Errorf("⚠️ API не вернул успешный ответ. HTTP: %s. Body: %s. Url: %s", resp.Status(), resp.Body(), resp.Request.URL)
	}

	return result, nil
}
