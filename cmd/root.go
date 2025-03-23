package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/itcaat/blet/config"
	tpclient "github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/cache"
	"github.com/joho/godotenv"
)

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

	if err := cache.Init(); err != nil {
		fmt.Println("❌ Ошибка инициализации кэша:", err)
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

		for _, city := range cache.CitiesCache {
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
		tpclient.GetCheapest(cfg.DefaultOrigin, token)

	case "week":
		fmt.Println("📅 Дешевые авиабилеты на неделю:")
		tpclient.GetWeekPrices(cfg.DefaultOrigin, token)

	default:
		fmt.Println("Неизвестный выбор")
	}
}
