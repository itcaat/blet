package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/cache"
	"github.com/itcaat/blet/internal/usecase"
)

func RunCheapest(cfg *config.Config, token string) {
	tickets, err := usecase.GetCheapestTickets(cfg.DefaultOrigin, cfg.DefaultDestination, strconv.FormatBool(cfg.OneWay), token)
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

	prepareTickets := func() {
		for _, t := range tickets {
			from := cache.GetCityName(t.Origin)
			to := cache.GetAnyName(t.Destination)
			route := fmt.Sprintf("%s → %s", from, to)
			if !cfg.OneWay {
				route += fmt.Sprintf(" → %s", from)
			}

			url, err := usecase.GetShortUrl(t.URL(), token)
			if err != nil {
				fmt.Println("❌ Ошибка:", err)
				return
			}
			partnerUrl := url[0].PartnerUrl

			desc := fmt.Sprintf("Туда: %s", t.DepartureAt)

			if cfg.OneWay {
				desc += fmt.Sprintf("— %d₽ — %s", t.Price, partnerUrl)
			} else {
				desc += fmt.Sprintf(". Обратно: %s — %d₽ — %s", t.ReturnAt, t.Price, partnerUrl)
			}

			if grouped[route] == nil {
				grouped[route] = []string{}
			}
			grouped[route] = append(grouped[route], desc)

			if details[route] == nil {
				details[route] = make(map[string]string)
			}
			details[route][desc] = t.URL()
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

// mapsKeys возвращает отсортированные ключи карты
func mapsKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// можно отсортировать, если хочешь алфавитный порядок
	return keys
}
