package analytics

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateClickEvent(ctx context.Context, event CreateClickEventParams) error
}

type CreateClickEventParams struct {
	LinkID    string
	ClickedAt time.Time
	Referrer  *string
	UserAgent *string
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateClickEvent(ctx context.Context, event CreateClickEventParams) error {
	if r.db == nil {
		return fmt.Errorf("create click event: db is nil")
	}

	clickedAt := event.ClickedAt.UTC()
	if clickedAt.IsZero() {
		clickedAt = time.Now().UTC()
	}

	model := ClickEvent{
		LinkID:    event.LinkID,
		ClickedAt: clickedAt,
		Referrer:  event.Referrer,
		UserAgent: event.UserAgent,
		CreatedAt: time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create click event: %w", err)
	}

	return nil
}
