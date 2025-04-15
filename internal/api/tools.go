package api

import (
	"fmt"

	"github.com/itcaat/blet/internal/models"
)

func (c *Client) GetShortUrlArray(tickets []*models.Ticket) error {
	const (
		apiURL    = "https://api.travelpayouts.com/links/v1/create"
		batchSize = 10
	)

	if len(tickets) == 0 {
		return nil
	}

	for start := 0; start < len(tickets); start += batchSize {
		end := start + batchSize
		if end > len(tickets) {
			end = len(tickets)
		}
		batch := tickets[start:end]

		links := make([]map[string]interface{}, len(batch))
		for i, t := range batch {
			links[i] = map[string]interface{}{"url": t.URL()}
		}

		var result models.ShortLinksResponse
		payload := map[string]interface{}{
			"trs":     400658,
			"marker":  616825,
			"shorten": true,
			"links":   links,
		}

		resp, err := c.resty.R().
			SetBody(payload).
			SetResult(&result).
			Post(apiURL)
		if err != nil {
			return fmt.Errorf("запрос не удался: %w", err)
		}

		if result.Status != "success" {
			return fmt.Errorf(
				"⚠️ API не вернул успешный ответ. HTTP: %s. Body: %s. URL: %s",
				resp.Status(), resp.Body(), resp.Request.URL,
			)
		}

		for i, t := range batch {
			if i < len(result.Result.Links) {
				t.ShortUrl = result.Result.Links[i].PartnerUrl
			}
		}
	}

	return nil
}
