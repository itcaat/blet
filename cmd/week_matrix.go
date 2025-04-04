package cmd

import (
	"fmt"

	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/cache"
	"github.com/itcaat/blet/internal/form"
	"github.com/itcaat/blet/internal/usecase"
)

func RunWeekMatrix(client *api.Client, cfg *config.Config) {

	var departDate, backDate string
	err := form.AskDates(&departDate, &backDate)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}
	fmt.Println("✈️ Дешевые авиабилеты на неделю:")
	fmt.Printf("Вылет-прилет: %s - %s\n", departDate, backDate)

	flights, err := usecase.GetWeekMatrix(client, cfg.DefaultOrigin, cfg.DefaultDestination, departDate, backDate)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	for _, flight := range flights {
		fmt.Printf("- %s → %s → %s за %d₽ (%s → %s, пересадок: %d)\n",
			cache.GetCityName(cfg.DefaultOrigin), cache.GetAnyName(flight.Destination), cache.GetCityName(cfg.DefaultOrigin), flight.Value,
			flight.DepartDate, flight.ReturnDate, flight.NumberOfStops)
	}
}
