package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr     string
	POSTGRES_DSN string
	REDIS_DSN    string

	PostgresPingTimeout time.Duration
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

	var err error
	cfg.PostgresPingTimeout, err = getenvDuration("POSTGRES_PING_TIMEOUT", 3*time.Second)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getenv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getenvDuration(key string, defaultVal time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultVal, nil
	}

	v, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be duration: %w", key, err)
	}
	return v, nil
}
