package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/form"
	"github.com/itcaat/blet/internal/models"
	"github.com/itcaat/blet/internal/usecase"
)

func RunWeekPrices(cfg *config.Config, token string) {
	dest := askDestination()
	var flights []models.WeekMatrixFlight
	var err error

	_ = spinner.New().
		Title("🔍 Ищем билеты на неделю...").
		Action(func() {
			flights, err = usecase.GetWeekMatrix(cfg.DefaultOrigin, dest, token)
			time.Sleep(500 * time.Millisecond)
		}).
		Run()

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	for _, flight := range flights {
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
