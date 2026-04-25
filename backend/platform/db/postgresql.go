package db

import (
	"fmt"
	"tracklink/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSQL(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.POSTGRES_DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgresql: %w", err)
	}

	return db, nil
}
