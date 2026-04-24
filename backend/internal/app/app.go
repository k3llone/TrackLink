package app

import (
	"fmt"

	"tracklink/internal/config"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	_ = cfg
	return nil
}
