package links

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, link *Link) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, link *Link) error {
	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}
