package cmd

import (
	"fmt"

	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/cache"
	"github.com/itcaat/blet/internal/usecase"
)

func RunSpecialOffers(client api.TravelpayoutsAPI, cfg *config.Config) {
	flights, err := usecase.GetSpecialOffers(client, cfg.DefaultOrigin)

	if err != nil {
		fmt.Println("❌ Ошибка:", err)
		return
	}

	for _, flight := range flights {
		resp, err := client.GetShortUrl(flight.URL())
		if err != nil {
			fmt.Println("❌ Ошибка:", err)
			return
		}
		partnerUrl := resp.Result.Links[0].PartnerUrl
		fmt.Printf("- %s → %s за %d₽ (Вылет: %s) %s\n",
			cache.GetCityName(cfg.DefaultOrigin), cache.GetAnyName(flight.Destination), flight.Price,
			flight.DepartDate, partnerUrl)
	}
}
