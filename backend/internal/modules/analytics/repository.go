package analytics

import (
	"context"
	"fmt"
	"time"

	"tracklink/internal/modules/links"

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

func (r *GormRepository) SumTotalClicks(ctx context.Context, ownerID string) (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("sum total clicks: db is nil")
	}

	type result struct {
		Value int64
	}
	var row result
	if err := r.db.WithContext(ctx).
		Model(&links.Link{}).
		Select("COALESCE(SUM(total_clicks), 0) AS value").
		Where("owner_id = ? AND status <> ? AND deleted_at IS NULL", ownerID, links.StatusDeleted).
		Scan(&row).Error; err != nil {
		return 0, fmt.Errorf("sum total clicks: %w", err)
	}

	return row.Value, nil
}

func (r *GormRepository) ListRecentLinks(ctx context.Context, ownerID string, limit int) ([]links.Link, error) {
	if r.db == nil {
		return nil, fmt.Errorf("list recent links: db is nil")
	}

	if limit <= 0 {
		limit = defaultRecentLinksLimit
	}

	items := make([]links.Link, 0, limit)
	if err := r.db.WithContext(ctx).
		Model(&links.Link{}).
		Where("owner_id = ? AND status <> ? AND deleted_at IS NULL", ownerID, links.StatusDeleted).
		Order("created_at DESC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list recent links: %w", err)
	}

	return items, nil
}
