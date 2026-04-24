package app

import (
	"fmt"
	"log"
	"net/http"

	"tracklink/internal/config"
	"tracklink/internal/httpapi"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	r := httpapi.NewRouter()
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
