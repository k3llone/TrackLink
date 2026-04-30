package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr      string
	PublicURL     string
	DatabaseURL   string
	RedisAddr     string
	SessionSecret string

	PostgresPingTimeout time.Duration
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	cfg := Config{
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		PublicURL:     getenv("PUBLIC_URL", "http://localhost:8080"),
		DatabaseURL:   getenvAny([]string{"DATABASE_URL", "POSTGRES_DSN"}, "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		RedisAddr:     normalizeRedisAddr(getenvAny([]string{"REDIS_ADDR", "REDIS_DSN"}, "localhost:6379")),
		SessionSecret: getenv("SESSION_SECRET", "dev-session-secret"),
	}

	var err error
	cfg.PostgresPingTimeout, err = getenvDuration("POSTGRES_PING_TIMEOUT", 3*time.Second)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getenvAny(keys []string, defaultVal string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return defaultVal
}

func getenv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func normalizeRedisAddr(raw string) string {
	if raw == "" {
		return "localhost:6379"
	}
	if strings.Contains(raw, "://") {
		return raw
	}
	if _, _, err := net.SplitHostPort(raw); err == nil {
		return "redis://" + raw
	}
	return raw
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
