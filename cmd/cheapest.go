package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/cache"
	"github.com/itcaat/blet/internal/usecase"
)

func RunCheapest(client api.TravelpayoutsAPI, cfg *config.Config) {
	tickets, err := usecase.GetCheapestTickets(client, cfg.DefaultOrigin, cfg.DefaultDestination, strconv.FormatBool(cfg.OneWay))
	if err != nil {
		fmt.Println("❌ Ошибка при получении данных:", err)
		return
	}

	if len(tickets) == 0 {
		fmt.Println("⚠️ Билеты не найдены.")
		return
	}

	grouped := make(map[string][]string)
	details := make(map[string]map[string]string)

	prepareTickets := func() {

		for _, t := range tickets {
			t := t
			from := cache.GetCityName(t.Origin)
			to := cache.GetAnyName(t.Destination)
			route := fmt.Sprintf("%s → %s", from, to)
			if !cfg.OneWay {
				route += fmt.Sprintf(" → %s", from)
			}

			desc := fmt.Sprintf("Туда: %s", t.DepartureAt)
			if cfg.OneWay {
				desc += fmt.Sprintf(" — %d₽", t.Price)
			} else {
				desc += fmt.Sprintf(". Обратно: %s — %d₽", t.ReturnAt, t.Price)
			}

			fullDesc := fmt.Sprintf("%s — %s", desc, t.ShortUrl)
			grouped[route] = append(grouped[route], fullDesc)
			if details[route] == nil {
				details[route] = make(map[string]string)
			}
			details[route][fullDesc] = t.URL()
		}
	}

	var selectedRoute string
	var selectedDesc string

	_ = spinner.New().Title("Ищем лучшие билетики...").Action(prepareTickets).Run()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("\n✈️ Самые дешевые авиабилеты").
				Description("Возвращает самые дешевые авиабилеты за определённые даты, найденные пользователями Авиасейлс за последние 48 часов."),
			huh.NewSelect[string]().
				Title("Выберите маршрут").
				Options(huh.NewOptions(mapsKeys(grouped)...)...).
				Height(5).
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
}

func mapsKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
