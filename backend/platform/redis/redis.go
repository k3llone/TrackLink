package redis

import (
	"context"
	"fmt"
	"time"
	"tracklink/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg config.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.REDIS_DSN)

	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	client := redis.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	return client, nil
}
