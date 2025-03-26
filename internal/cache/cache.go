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

// --- Типы данных ---

type City struct {
	Name                 string `json:"name"`
	Code                 string `json:"code"`
	CountryCode          string `json:"country_code"`
	HasFlightableAirport bool   `json:"has_flightable_airport"`
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// --- Обобщённый кеш ---

type genericCache[T any] struct {
	URL      string
	FileName string
	Data     []T
	once     sync.Once
	mu       sync.RWMutex
}

// --- Интерфейс для инициализации и реестр всех кешей ---

type initializableCache interface {
	init() error
}

var allCaches []initializableCache

func newGenericCache[T any](url, fileName string) *genericCache[T] {
	cache := &genericCache[T]{
		URL:      url,
		FileName: fileName,
	}
	allCaches = append(allCaches, cache)
	return cache
}

// --- Конкретные кеши ---

var (
	citiesCache    = newGenericCache[City]("https://api.travelpayouts.com/data/ru/cities.json", "cities.json")
	countriesCache = newGenericCache[Country]("https://api.travelpayouts.com/data/ru/countries.json", "countries.json")
)

// --- Публичная инициализация всех кешей ---

func Init() error {
	for _, c := range allCaches {
		if err := c.init(); err != nil {
			return err
		}
	}
	return nil
}

// --- Внутренняя инициализация конкретного кеша ---

func (c *genericCache[T]) init() error {
	var err error
	c.once.Do(func() {
		var path string
		path, err = ensureCacheFile(c.URL, c.FileName)
		if err != nil {
			return
		}

		var data []byte
		data, err = os.ReadFile(path)
		if err != nil {
			return
		}

		var parsed []T
		if err = json.Unmarshal(data, &parsed); err != nil {
			return
		}

		c.mu.Lock()
		defer c.mu.Unlock()
		c.Data = parsed
	})
	return err
}

// --- Загрузка и обновление файла при необходимости ---

func ensureCacheFile(url, fileName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(home, ".blet", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(cacheDir, fileName)

	stat, err := os.Stat(fullPath)
	if err == nil && stat.ModTime().AddDate(0, 0, 7).After(time.Now()) {
		return fullPath, nil
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки %s: %w", fileName, err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("неудачный ответ при загрузке %s: %s", fileName, resp.Status())
	}

	if err := os.WriteFile(fullPath, resp.Body(), 0644); err != nil {
		return "", fmt.Errorf("не удалось сохранить %s: %w", fileName, err)
	}

	return fullPath, nil
}

// --- Публичные функции доступа ---

func GetCityName(code string) string {
	citiesCache.mu.RLock()
	defer citiesCache.mu.RUnlock()

	for _, city := range citiesCache.Data {
		if city.Code == code {
			return city.Name
		}
	}
	return code
}

func GetCountryName(code string) string {
	countriesCache.mu.RLock()
	defer countriesCache.mu.RUnlock()

	for _, country := range countriesCache.Data {
		if country.Code == code {
			return country.Name
		}
	}
	return code
}

func (c *genericCache[T]) Filter(predicate func(T) bool) []T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []T
	for _, item := range c.Data {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Cities() *genericCache[City] {
	return citiesCache
}

func Countries() *genericCache[Country] {
	return countriesCache
}
