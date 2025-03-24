package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/cache"
	cmdhelpers "github.com/itcaat/blet/internal/cmd"
	"github.com/itcaat/blet/internal/usecase"
)

func RunCheapest(cfg *config.Config, token string) {
	tickets, err := usecase.GetCheapestTickets(cfg.DefaultOrigin, token)
	if err != nil {
		fmt.Println("❌ Ошибка при получении данных:", err)
		return
	}

	if len(tickets) == 0 {
		fmt.Println("⚠️ Билеты не найдены.")
		return
	}

	// Группировка по маршрутам
	grouped := make(map[string][]string)          // "Москва → Сочи" -> список описаний
	details := make(map[string]map[string]string) // [маршрут][описание] -> ссылка

	for _, t := range tickets {
		from := cache.GetHumanCityName(t.Origin)
		to := cache.GetHumanCityName(t.Destination)
		route := fmt.Sprintf("%s → %s", from, to)
		desc := fmt.Sprintf("%s — %d₽", t.DepartureAt, t.Price)

		if grouped[route] == nil {
			grouped[route] = []string{}
		}
		grouped[route] = append(grouped[route], desc)

		if details[route] == nil {
			details[route] = make(map[string]string)
		}
		details[route][desc] = cmdhelpers.FormatAviasalesLink(t.Link)
	}

	var selectedRoute string
	var selectedDesc string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Выберите маршрут").
				Options(huh.NewOptions(mapsKeys(grouped)...)...).
				Height(8).
				Value(&selectedRoute),
			huh.NewSelect[string]().
				Title("Выберите рейс").
				Height(8).
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(grouped[selectedRoute]...)
				}, &selectedRoute).
				Value(&selectedDesc),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Println("❌ Ошибка:", err)
		os.Exit(1)
	}

	link := details[selectedRoute][selectedDesc]
	fmt.Printf("\n🔗 Ссылка на билет: %s\n", link)
}

// mapsKeys возвращает отсортированные ключи карты
func mapsKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// можно отсортировать, если хочешь алфавитный порядок
	return keys
}
