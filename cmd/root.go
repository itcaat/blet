package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/itcaat/blet/config"
	"github.com/itcaat/blet/internal/api"
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
	if err != nil || cfg.DefaultOrigin == "" || cfg.DefaultDestination == "" {
		cfg.DefaultOrigin = "MOW"
		cfg.DefaultDestination = "LED"
		cfg.OneWay = true

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("❌ Не удалось сохранить конфиг:", err)
			os.Exit(1)
		}
	}

	// форма выбора города вылета

	var change_default_origin bool

	form_change_default_origin := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("\nПривествую, странник. Кажется, пора полетать!? ✈️"),
			huh.NewConfirm().
				Title(fmt.Sprintf("Откуда: %s \nКуда: %s\nНаправление: %s \n\nОставим как есть или поменяем?", cache.GetCityName(cfg.DefaultOrigin), cache.GetAnyName(cfg.DefaultDestination), cfg.TekstomPokazhiOneWay())).
				Value(&change_default_origin).
				Affirmative("Выбрать другой").
				Negative("Оставить"),
		))

	if err := form_change_default_origin.Run(); err != nil {
		log.Fatal(err)
	}

	if change_default_origin {
		form.ShowIataTargets(&cfg.DefaultOrigin, "Откуда полетим", false)
		form.ShowIataTargets(&cfg.DefaultDestination, "Куда полетим (можно выбрать страну или город)", true)
		form.ShowConfirm(&cfg.OneWay, "В одну с сторону?", "Только туда", "Туда-обратно")

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("❌ Не удалось сохранить конфиг:", err)
			os.Exit(1)
		}
	}

	fmt.Sprintln(cfg.OneWay)

	// emoji airplane
	var choice string

	form := huh.NewForm(

		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("%s ➡️  %s", cache.GetCityName(cfg.DefaultOrigin), cache.GetAnyName(cfg.DefaultDestination))).
				Options(
					huh.NewOption("✈️ Самые дешевые авиабилеты", "cheapest"),
					huh.NewOption("👻 Поиск по недельной матрице", "week"),
					huh.NewOption("🔮 Спецпредложения", "special"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	client := api.New(token)

	switch choice {
	case "cheapest":
		RunCheapest(client, &cfg)

	case "week":
		RunWeekMatrix(client, &cfg)

	case "special":
		fmt.Println("✈️ Спецпредложения от авиакомпаний:")
		RunSpecialOffers(client, &cfg)

	default:
		fmt.Println("Неизвестный выбор")
	}
}
