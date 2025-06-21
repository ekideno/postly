package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database Database `yaml:"database"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	User     string `yaml:"user" env-required:"true"`
	Password int    `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
	Sslmode  string `yaml:"sslmode" env-required:"disable"`
	Timezone string `yaml:"timezone" env-default:"UTC"`
}

func New() *Config {
	// Get path of config file
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		log.Fatalf("CONFIG_PATH env value not found")
	}

	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Failed to find config file: %s", path)
	}

	// Read config file
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	return &cfg
}
