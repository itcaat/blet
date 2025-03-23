package tpclient

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

func GetCheapest(origin, token string) {

	var result models.PriceForDatesResponse

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origin":   origin,
			"one_way":  "true",
			"currency": "rub",
			"limit":    "5",
			"token":    token,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get(pricesForDatesURL)

	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
		return
	}

	if !result.Success {
		fmt.Printf("⚠️ Неуспешный ответ API. Статус: %s\n", resp.Status())
		return
	}

	for _, t := range result.Data {
		fmt.Printf("- %s → %s за %d₽ (%s)\n", t.Origin, t.Destination, t.Price, t.DepartureAt)
		fmt.Printf("  Ссылка: %s\n", formatAviasalesLink(t.Link))
	}
}

func GetWeekPrices(origin, token string) {
	var result models.WeekMatrixResponse

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origin":      origin,
			"destination": "LED",
			"currency":    "rub",
			"depart_date": "2025-09-04",
			"return_date": "2025-09-11",
			"token":       token,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get(weekMatrixURL)

	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
		return
	}

	if !result.Success {
		fmt.Printf("⚠️ Неуспешный ответ API. Статус: %s. Ответ: %s\n", resp.Status(), resp.Body())
		return
	}

	for _, flight := range result.Data {
		fmt.Printf("- %s → %s за %d₽ (%s → %s, пересадок: %d)\n",
			origin, flight.Destination, flight.Value,
			flight.DepartDate, flight.ReturnDate, flight.NumberOfStops)
	}
}

func formatAviasalesLink(path string) string {
	return fmt.Sprintf("https://www.aviasales.ru%s", path)
}
