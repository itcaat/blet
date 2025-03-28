package form

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/itcaat/blet/internal/cache"
)

func ShowIataTargets(selectedIATA *string, titleSelect string, includeCountry bool) {
	var cityPairs []struct {
		Label string
		Code  string
	}

	for _, city := range cache.Cities().Filter(func(c cache.City) bool {
		return c.HasFlightableAirport
	}) {
		label := fmt.Sprintf("%s (%s)", city.Name, city.Code)
		cityPairs = append(cityPairs, struct {
			Label string
			Code  string
		}{Label: label, Code: city.Code})
	}

	if includeCountry {
		for _, country := range cache.Countries().Data {
			label := fmt.Sprintf("%s (%s)", country.Name, country.Code)
			cityPairs = append(cityPairs, struct {
				Label string
				Code  string
			}{Label: label, Code: country.Code})
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

	//надо через временную переменную, иначе есть баг
	var selected string

	// UI выбора города
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(titleSelect).
				Height(10).
				Options(options...).
				Value(&selected),
		),
	)
	if err := form.Run(); err != nil {
		fmt.Println("❌ Ошибка выбора:", err)
		os.Exit(1)
	}
	*selectedIATA = selected

}
