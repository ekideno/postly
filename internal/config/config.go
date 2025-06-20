package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// TODO: Move all database values to nested struct as the main struct expands
type Config struct {
	DB_HOST     string `yaml:"env" env-default:"localhost"`
	DB_USER     string `yaml:"user" env-required:"true"`
	DB_PASSWORD int    `yaml:"password" env-required:"true"`
	DB_NAME     string `yaml:"name" env-required:"true"`
	DB_PORT     int    `yaml:"port" env-default:"5432"`
	DB_SSLMODE  string `yaml:"sslmode" env-required:"disable"`
	DB_TIMEZONE string `yaml:"timezone" env-default:"UTC"`
}

func New() *Config {
	// Get path of config file
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("CONFIG_PATH env value not found")
	}

	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("Failed to find config file: %s", path))
	}

	// Read config file
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(fmt.Sprintf("Failed to read config file: %s", err.Error()))
	}

	return &cfg
}
