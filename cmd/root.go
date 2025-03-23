package cmd

import (
	"fmt"
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

		fmt.Println("✅ Город вылета установлен:", cfg.DefaultOrigin)
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
		RunCheapest(&cfg, token)

	case "week":
		fmt.Println("📅 Дешевые авиабилеты на неделю:")
		RunWeekPrices(&cfg, token)

	default:
		fmt.Println("Неизвестный выбор")
	}
}
