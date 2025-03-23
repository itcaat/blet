package form

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/itcaat/blet/internal/cache"
)

func ShowCityPairs(selectedIATA *string, titleSelect string) {
	var cityPairs []struct {
		Label string
		Code  string
	}

	for _, city := range cache.CitiesCache {
		if city.HasFlightableAirport {
			label := fmt.Sprintf("%s (%s)", city.Name, city.Code)
			cityPairs = append(cityPairs, struct {
				Label string
				Code  string
			}{Label: label, Code: city.Code})
		}
	}

	// Сортируем по названию
	sort.Slice(cityPairs, func(i, j int) bool {
		return cityPairs[i].Label < cityPairs[j].Label
	})

	// Преобразуем в huh.Options
	var options []huh.Option[string]
	for _, pair := range cityPairs {
		options = append(options, huh.NewOption(pair.Label, pair.Code))
	}

	// UI выбора города
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(titleSelect),
			huh.NewSelect[string]().
				Options(options...).
				Value(selectedIATA),
		),
	)
	if err := form.Run(); err != nil {
		fmt.Println("❌ Ошибка выбора:", err)
		os.Exit(1)
	}

}
