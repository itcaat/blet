package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type City struct {
	Name                 string `json:"name"`
	Code                 string `json:"code"`
	CountryCode          string `json:"country_code"`
	HasFlightableAirport bool   `json:"has_flightable_airport"`
}

var (
	CitiesCache []City
	once        sync.Once
)

const citiesURL = "https://api.travelpayouts.com/data/ru/cities.json"

// Init проверяет наличие cities.json и загружает его в память
func Init() error {
	var err error
	once.Do(func() {
		var path string
		path, err = ensureCitiesFile()
		if err != nil {
			return
		}

		var data []byte
		data, err = os.ReadFile(path)
		if err != nil {
			return
		}

		err = json.Unmarshal(data, &CitiesCache)
	})
	return err
}

func ensureCitiesFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(home, ".blet", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	citiesPath := filepath.Join(cacheDir, "cities.json")

	// Проверим наличие и актуальность (например, возраст файла > 7 дней)
	stat, err := os.Stat(citiesPath)
	if err == nil {
		if stat.ModTime().AddDate(0, 0, 7).After(time.Now()) {
			return citiesPath, nil // файл свежий
		}
	}

	// Загружаем файл
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(citiesURL)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки cities.json: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("неудачный ответ при загрузке cities.json: %s", resp.Status())
	}

	if err := os.WriteFile(citiesPath, resp.Body(), 0644); err != nil {
		return "", fmt.Errorf("не удалось сохранить cities.json: %w", err)
	}

	return citiesPath, nil
}

func GetHumanCityName(code string) string {
	for _, city := range CitiesCache {
		if city.Code == code {
			return city.Name
		}
	}
	return code
}
