package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CONFIG_FILE = "blog.json"
)

type Config struct {
	Version     int      `json:"version"`
	Domain      string   `json:"domain"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ChillFiles  []string `json:"chill-files"`
	Footer      string   `json:"footer"`
	Discussion  bool     `json:"discussion"`
}

func CreateDefaultConfig() *Config {
	return &Config{
		Version:     1,
		Domain:      "{Your blogctl domain}",
		Author:      "{Your name}",
		Title:       "{Title of your blog}",
		Description: "{Description of your blog}",
		ChillFiles:  []string{},
		Footer:      "{Footer of your blog}",
		Discussion:  false,
	}
}

func ConfigPath(postPath string) string {
	return filepath.Join(postPath, CONFIG_FILE)
}

func Load(postPath string) (*Config, error) {
	raw, err := os.ReadFile(ConfigPath(postPath))
	if err != nil {
		return nil, fmt.Errorf("Error reading config: %s", err)
	}

	config := new(Config)
	if err := json.Unmarshal(raw, &config); err != nil {
		return nil, fmt.Errorf("Error parsing config: %s", err)
	}

	return config, nil
}
