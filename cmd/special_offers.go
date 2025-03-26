package cmd

import (
	"fmt"

	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/cache"
	"github.com/itcaat/blet/internal/usecase"
)

func RunSpecialOffers(cfg *config.Config, token string) {
	flights, err := usecase.GetSpecialOffers(cfg.DefaultOrigin, token)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	for _, flight := range flights {
		fmt.Printf("- %s → %s за %d₽ (Вылет: %s)\n",
			cache.GetHumanCityName(cfg.DefaultOrigin), cache.GetHumanCityName(flight.Destination), flight.Price,
			flight.DepartDate)
	}
}
