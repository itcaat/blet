package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/cache"
	"github.com/itcaat/blet/internal/form"
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
		form.ShowCityPairs(&cfg.DefaultOrigin, "Давай выберем город вылета по-умолчанию")
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("❌ Не удалось сохранить конфиг:", err)
			os.Exit(1)
		}

		fmt.Println("✅ IATA код города установлен: ", cfg.DefaultOrigin)
	}

	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("👋 Какие билеты будем искать? Город вылета по умолчанию: %s", cache.GetHumanCityName(cfg.DefaultOrigin))).
				Options(
					huh.NewOption("Билеты хоть куда", "cheapest"),
					huh.NewOption("Поиск по недельной матрице", "week"),
					huh.NewOption("Спецпредложения", "special"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	switch choice {
	case "cheapest":
		fmt.Println("✈️ Самые дешевые авиабилеты:")
		RunCheapest(&cfg, token)

	case "week":
		fmt.Println("📅 Дешевые авиабилеты на неделю:")
		RunWeekPrices(&cfg, token)

	case "special":
		fmt.Println("📅 Дешевые авиабилеты на неделю:")
		RunSpecialOffers(&cfg, token)

	default:
		fmt.Println("Неизвестный выбор")
	}
}
