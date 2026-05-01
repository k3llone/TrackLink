package links

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, link *Link) error
	ExistsByCode(ctx context.Context, code string) (bool, error)
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

func (r *GormRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var link Link
	err := r.db.WithContext(ctx).
		Select("id").
		Where("code = ?", code).
		First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("check link code existence: %w", err)
	}

	return true, nil
}
