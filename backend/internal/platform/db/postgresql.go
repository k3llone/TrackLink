package db

import (
	"context"
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

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("postgresql sql db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PostgresPingTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("postgresql ping: %w", err)
	}

	return db, nil
}
