package config

import (
	"os"
	"time"
)

type Config struct {
	Port           string
	TimeoutDuration time.Duration
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:           port,
		TimeoutDuration: 30 * time.Second,
	}
}