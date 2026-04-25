package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr     string
	POSTGRES_DSN string
	REDIS_DSN    string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	cfg := Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		POSTGRES_DSN: getenv("POSTGRES_DSN", "postgresql://postgres:postgres@localhost:5432/postgres"),
		REDIS_DSN:    getenv("REDIS_DSN", "redis://localhost:6379"),
	}

	return cfg, nil
}

func getenv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
