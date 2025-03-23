package cmd

import (
	"fmt"

	"github.com/itcaat/blet/config"
	tpclient "github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/form"
)

func RunWeekPrices(cfg *config.Config, token string) {
	destination := askDestination()

	result, err := tpclient.GetWeekPrices(cfg.DefaultOrigin, destination, token)
	if err != nil {
		fmt.Println("❌ Ошибка при получении данных:", err)
		return
	}

	for _, flight := range result.Data {
		fmt.Printf("- %s → %s за %d₽ (%s → %s, пересадок: %d)\n",
			cfg.DefaultOrigin, flight.Destination, flight.Value,
			flight.DepartDate, flight.ReturnDate, flight.NumberOfStops)
	}
}

func askDestination() string {
	var dest string
	form.ShowCityPairs(&dest, "Куда летим?")
	return dest
}
