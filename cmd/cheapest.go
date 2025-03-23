package cmd

import (
	"fmt"

	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/cmd"
	"github.com/itcaat/blet/internal/usecase"
)

func RunCheapest(cfg *config.Config, token string) {
	tickets, err := usecase.GetCheapestTickets(cfg.DefaultOrigin, token)
	if err != nil {
		fmt.Println("❌ Ошибка при получении данных:", err)
		return
	}

	for _, t := range tickets {
		fmt.Printf("- %s → %s за %d₽ (%s)\n", t.Origin, t.Destination, t.Price, t.DepartureAt)
		fmt.Printf("  Ссылка: %s\n", cmd.FormatAviasalesLink(t.Link))
	}
}
