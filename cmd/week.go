package cmd

import (
	"fmt"

	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/cache"
	"github.com/itcaat/blet/internal/form"
	"github.com/itcaat/blet/internal/usecase"
)

func RunWeekPrices(cfg *config.Config, token string) {
	dest := askDestination()
	depart, back, err := form.AskDates()

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	fmt.Printf("Вылет-прилет: %s - %s\n", depart, back)

	flights, err := usecase.GetWeekMatrix(cfg.DefaultOrigin, dest, depart, back, token)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	for _, flight := range flights {
		fmt.Printf("- %s → %s → %s за %d₽ (%s → %s, пересадок: %d)\n",
			cache.GetHumanCityName(cfg.DefaultOrigin), cache.GetHumanCityName(flight.Destination), cache.GetHumanCityName(cfg.DefaultOrigin), flight.Value,
			flight.DepartDate, flight.ReturnDate, flight.NumberOfStops)
	}
}

func askDestination() string {
	var dest string
	form.ShowCityPairs(&dest, "Куда летим?")
	return dest
}
