package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/go-resty/resty/v2"
	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/cache"
	"github.com/joho/godotenv"
)

type Ticket struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Price       int    `json:"price"`
	DepartureAt string `json:"departure_at"`
	Link        string `json:"link"`
}

type PriceForDatesResponse struct {
	Success  bool     `json:"success"`
	Data     []Ticket `json:"data"`
	Currency string   `json:"currency"`
}

func Execute() {
	if len(os.Args) > 1 && os.Args[1] == "--reset" {
		home, _ := os.UserHomeDir()
		bletPath := filepath.Join(home, ".blet")

		if err := os.RemoveAll(bletPath); err != nil {
			fmt.Println("❌ Не удалось удалить ~/.blet:", err)
			os.Exit(1)
		}

		fmt.Println("🧹 Конфигурация сброшена. Папка ~/.blet удалена.")
		os.Exit(0)
	}

	citiesPath, err := cache.EnsureCitiesCache()
	if err != nil {
		fmt.Println("❌ cities.json error:", err)
		os.Exit(1)
	}

	cities, err := cache.LoadCities(citiesPath)
	if err != nil {
		fmt.Println("❌ Ошибка парсинга cities.json:", err)
		os.Exit(1)
	}

	// Загружаем .env
	_ = godotenv.Load()
	token := os.Getenv("AVIASALES_TOKEN")
	if token == "" {
		fmt.Println("❌ Переменная AVIASALES_TOKEN не задана в .env")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig()
	if err != nil || cfg.DefaultOrigin == "" {
		// Готовим список (название + код)
		var selectedIATA string
		var cityPairs []struct {
			Label string
			Code  string
		}

		for _, city := range cities {
			if city.HasFlightableAirport && city.CountryCode == "RU" {
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
				huh.NewSelect[string]().
					Title("Давай выберем город вылета по-умолчанию").
					Options(options...).
					Value(&selectedIATA),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Println("❌ Ошибка выбора:", err)
			os.Exit(1)
		}

		// Сохраняем конфиг
		cfg.DefaultOrigin = selectedIATA
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("❌ Не удалось сохранить конфиг:", err)
			os.Exit(1)
		}

		fmt.Println("✅ Город вылета установлен:", selectedIATA)
	} else {
		fmt.Printf("🌍 Город вылета по умолчанию: %s\n", cfg.DefaultOrigin)
	}

	var choice string
	menu := huh.NewSelect[string]().
		Title("Выберите действие").
		Options(
			huh.NewOption("Самые дешевые авиабилеты", "cheapest"),
			huh.NewOption("Дешевые авиабилеты на неделю", "week"),
		).
		Value(&choice)

	if err := menu.Run(); err != nil {
		fmt.Println("Ошибка выбора:", err)
		os.Exit(1)
	}

	switch choice {
	case "cheapest":
		fmt.Println("✈️ Самые дешевые авиабилеты:")
		getCheapest(cfg.DefaultOrigin, token)

	case "week":
		fmt.Println("📅 Дешевые авиабилеты на неделю:")
		getWeekPrices(cfg.DefaultOrigin, token)

	default:
		fmt.Println("Неизвестный выбор")
	}
}

func getCheapest(origin, token string) {
	client := resty.New()

	var result PriceForDatesResponse

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origin":   origin,
			"one_way":  "true",
			"currency": "rub",
			"limit":    "5",
			"token":    token,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get("https://api.travelpayouts.com/aviasales/v3/prices_for_dates")

	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
		return
	}

	if !result.Success {
		fmt.Printf("⚠️ Неуспешный ответ API. Статус: %s\n", resp.Status())
		return
	}

	for _, t := range result.Data {
		fmt.Printf("- %s → %s за %d₽ (%s)\n", t.Origin, t.Destination, t.Price, t.DepartureAt)
		fmt.Printf("  Ссылка: https://www.aviasales.ru%s\n", t.Link)
	}
}

func getWeekPrices(origin, token string) {
	client := resty.New()

	var result struct {
		Success bool `json:"success"`
		Data    []struct {
			Destination   string `json:"destination"`
			DepartDate    string `json:"depart_date"`
			ReturnDate    string `json:"return_date"`
			Value         int    `json:"value"`
			NumberOfStops int    `json:"number_of_changes"`
		} `json:"data"`
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"origin":             origin,
			"destination":        "LED",
			"currency":           "rub",
			"depart_date":        "2025-09-04",
			"return_date":        "2025-09-11",
			"show_to_affiliates": "true",
			"token":              token,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Get("https://api.travelpayouts.com/v2/prices/week-matrix")

	if err != nil {
		fmt.Println("Ошибка при запросе:", err)
		return
	}

	if !result.Success {
		fmt.Printf("⚠️ Неуспешный ответ API. Статус: %s\n", resp.Status())
		return
	}

	for _, flight := range result.Data {
		fmt.Printf("- %s → %s за %d₽ (%s → %s, пересадок: %d)\n",
			origin, flight.Destination, flight.Value,
			flight.DepartDate, flight.ReturnDate, flight.NumberOfStops)
	}
}
