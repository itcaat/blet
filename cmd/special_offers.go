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
		url, err := usecase.GetShortUrl(flight.URL(), token)
		if err != nil {
			fmt.Println("❌ Ошибка:", err)
			return
		}
		partnerUrl := url[0].PartnerUrl
		fmt.Printf("- %s → %s за %d₽ (Вылет: %s) %s\n",
			cache.GetCityName(cfg.DefaultOrigin), cache.GetCityName(flight.Destination), flight.Price,
			flight.DepartDate, partnerUrl)
	}
}
