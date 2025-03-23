package cmd

import (
	"fmt"

	"github.com/itcaat/blet/config"
	tpclient "github.com/itcaat/blet/internal/api"
	cmdhelpers "github.com/itcaat/blet/internal/cmd"
)

func RunCheapest(cfg *config.Config, token string) {
	result, err := tpclient.GetCheapest(cfg.DefaultOrigin, token)
	if err != nil {
		fmt.Println("❌ Ошибка при получении данных:", err)
		return
	}

	for _, t := range result.Data {
		fmt.Printf("- %s → %s за %d₽ (%s)\n", t.Origin, t.Destination, t.Price, t.DepartureAt)
		fmt.Printf("  Ссылка: %s\n", cmdhelpers.FormatAviasalesLink(t.Link))
	}
}
