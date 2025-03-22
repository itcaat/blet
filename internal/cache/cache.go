package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	citiesURL   = "https://api.travelpayouts.com/data/ru/cities.json"
	cacheMaxAge = 7 * 24 * time.Hour
)

// Структура города
type City struct {
	Name                 string `json:"name"`
	Code                 string `json:"code"`
	CountryCode          string `json:"country_code"`
	HasFlightableAirport bool   `json:"has_flightable_airport"`
}

func EnsureCitiesCache() (string, error) {
	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".blet", "cache")
	citiesPath := filepath.Join(cacheDir, "cities.json")

	if info, err := os.Stat(citiesPath); err == nil {
		if time.Since(info.ModTime()) < cacheMaxAge {
			fmt.Println("📍 Найден актуальный cities.json")
			return citiesPath, nil
		}
		fmt.Println("🔄 Кеш устарел — перекачиваем cities.json...")
	} else {
		fmt.Println("⬇️ Качаем cities.json...")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	resp, err := http.Get(citiesURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(citiesPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("✅ cities.json успешно скачан")
	return citiesPath, nil
}

func LoadCities(path string) ([]City, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cities []City
	if err := json.Unmarshal(data, &cities); err != nil {
		return nil, err
	}

	var filtered []City
	for _, c := range cities {
		if c.HasFlightableAirport {
			filtered = append(filtered, c)
		}
	}

	return filtered, nil
}
