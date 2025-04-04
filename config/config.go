package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultOrigin      string `json:"default_origin"`
	DefaultDestination string `json:"default_destination"`
	OneWay             bool   `json:"one_way"`
}

func (c *Config) TekstomPokazhiOneWay() string {
	if c.OneWay {
		return "В одну сторону"
	} else {
		return "Туда-обратно"
	}
}

func configPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".blet")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

func SaveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadConfig() (Config, error) {
	var cfg Config
	path, err := configPath()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
