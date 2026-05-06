package app

import (
	"fmt"
	"net/http"

	"tracklink/internal/config"
	httpserver "tracklink/internal/http"
	"tracklink/internal/platform/db"
	applogger "tracklink/internal/platform/logger"
	platformredis "tracklink/internal/platform/redis"
	"tracklink/internal/platform/session"
)

func Run() error {
	log := applogger.New("text")

	cfg, err := config.Load()
	if err != nil {
		log.Error("config_load_failed", "error", err)
		return fmt.Errorf("config: %w", err)
	}
	log = applogger.New(cfg.LogFormat)

	postgresDB, err := db.NewPostgreSQL(cfg)
	if err != nil {
		log.Error("postgres_init_failed", "error", err)
		return fmt.Errorf("postgres init: %w", err)
	}

	redisClient, err := platformredis.NewRedis(cfg)
	if err != nil {
		log.Error("redis_init_failed", "error", err)
		return fmt.Errorf("redis init: %w", err)
	}

	sessionStore := session.NewRedisStore(redisClient)

	r := httpserver.NewRouter(log, httpserver.Deps{
		DB:       postgresDB,
		Redis:    redisClient,
		Sessions: sessionStore,
		Config:   cfg,
	})
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}
	log.Info("app_start", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server_start_failed", "error", err)
		return fmt.Errorf("http server: %w", err)
	}
	log.Info("app_stop")
	return nil
}
