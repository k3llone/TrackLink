package app

import (
	"fmt"
	"log"
	"net/http"

	"tracklink/internal/config"
	httpserver "tracklink/internal/http"
	"tracklink/internal/platform/db"
	platformredis "tracklink/internal/platform/redis"
	"tracklink/internal/platform/session"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	postgresDB, err := db.NewPostgreSQL(cfg)
	if err != nil {
		return fmt.Errorf("postgres init: %w", err)
	}

	redisClient, err := platformredis.NewRedis(cfg)
	if err != nil {
		return fmt.Errorf("redis init: %w", err)
	}

	sessionStore := session.NewRedisStore(redisClient)

	r := httpserver.NewRouter(httpserver.Deps{
		DB:       postgresDB,
		Redis:    redisClient,
		Sessions: sessionStore,
	})
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}
	log.Printf("http listening on %s", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}
