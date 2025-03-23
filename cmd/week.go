package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh/spinner"
	"github.com/itcaat/blet/config"
	tpclient "github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/form"
	"github.com/itcaat/blet/internal/models"
)

func RunWeekPrices(cfg *config.Config, token string) {
	destination := askDestination()

	var result *models.WeekMatrixResponse
	var apiErr error

	action := func() {
		result, apiErr = tpclient.GetWeekPrices(cfg.DefaultOrigin, destination, token)
	}

	_ = spinner.New().
		Title("🔍 Ищем билеты на неделю...").
		Action(action).
		Run()

	if apiErr != nil {
		fmt.Println("❌ Ошибка при получении данных:", apiErr)
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
