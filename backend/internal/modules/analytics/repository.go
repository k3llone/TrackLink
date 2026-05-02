package analytics

import (
	"context"
	"errors"
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

func (r *GormRepository) CountTotalLinks(ctx context.Context, ownerID string) (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("count total links: db is nil")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&links.Link{}).
		Where("owner_id = ? AND status <> ? AND deleted_at IS NULL", ownerID, links.StatusDeleted).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count total links: %w", err)
	}

	return count, nil
}

func (r *GormRepository) CountActiveLinks(ctx context.Context, ownerID string) (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("count active links: db is nil")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&links.Link{}).
		Where("owner_id = ? AND status = ? AND deleted_at IS NULL", ownerID, links.StatusActive).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count active links: %w", err)
	}

	return count, nil
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

func (r *GormRepository) CountClicksSince(ctx context.Context, ownerID string, since time.Time) (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("count clicks since: db is nil")
	}

	var count int64
	err := r.db.WithContext(ctx).
		Table("click_events").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.owner_id = ? AND click_events.clicked_at >= ?", ownerID, since.UTC()).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count clicks since: %w", err)
	}

	return count, nil
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

func (r *GormRepository) GetLinkByIDAndOwner(ctx context.Context, linkID, ownerID string) (links.Link, error) {
	if r.db == nil {
		return links.Link{}, fmt.Errorf("get link by id and owner: db is nil")
	}

	var link links.Link
	err := r.db.WithContext(ctx).
		Model(&links.Link{}).
		Where("id = ? AND owner_id = ? AND status <> ? AND deleted_at IS NULL", linkID, ownerID, links.StatusDeleted).
		First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return links.Link{}, ErrLinkNotFound
		}
		return links.Link{}, fmt.Errorf("get link by id and owner: %w", err)
	}

	return link, nil
}

func (r *GormRepository) CountLinkClicks(ctx context.Context, linkID string, from, to time.Time) (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("count link clicks: db is nil")
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&ClickEvent{}).
		Where("link_id = ? AND clicked_at >= ? AND clicked_at <= ?", linkID, from.UTC(), to.UTC()).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count link clicks: %w", err)
	}

	return count, nil
}

func (r *GormRepository) CountLinkClicksSince(ctx context.Context, linkID string, since time.Time) (int64, error) {
	if r.db == nil {
		return 0, fmt.Errorf("count link clicks since: db is nil")
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&ClickEvent{}).
		Where("link_id = ? AND clicked_at >= ?", linkID, since.UTC()).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count link clicks since: %w", err)
	}

	return count, nil
}

func (r *GormRepository) LastLinkClickedAt(ctx context.Context, linkID string, from, to time.Time) (*time.Time, error) {
	if r.db == nil {
		return nil, fmt.Errorf("last link clicked at: db is nil")
	}

	type result struct {
		Value *time.Time
	}
	var row result
	err := r.db.WithContext(ctx).
		Model(&ClickEvent{}).
		Select("MAX(clicked_at) AS value").
		Where("link_id = ? AND clicked_at >= ? AND clicked_at <= ?", linkID, from.UTC(), to.UTC()).
		Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("last link clicked at: %w", err)
	}
	if row.Value == nil {
		return nil, nil
	}

	value := row.Value.UTC()
	return &value, nil
}

func (r *GormRepository) ListLinkClickSeries(ctx context.Context, linkID string, from, to time.Time, groupBy string) ([]TimeSeriesBucket, error) {
	if r.db == nil {
		return nil, fmt.Errorf("list link click series: db is nil")
	}

	truncUnit := GroupByDay
	if groupBy == GroupByHour {
		truncUnit = GroupByHour
	}

	type row struct {
		PeriodStart time.Time
		Clicks      int64
	}
	rows := make([]row, 0)
	err := r.db.WithContext(ctx).
		Model(&ClickEvent{}).
		Select(fmt.Sprintf("date_trunc('%s', clicked_at) AS period_start, COUNT(*) AS clicks", truncUnit)).
		Where("link_id = ? AND clicked_at >= ? AND clicked_at <= ?", linkID, from.UTC(), to.UTC()).
		Group("period_start").
		Order("period_start ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list link click series: %w", err)
	}

	buckets := make([]TimeSeriesBucket, 0, len(rows))
	for _, row := range rows {
		buckets = append(buckets, TimeSeriesBucket{
			PeriodStart: row.PeriodStart.UTC(),
			Clicks:      row.Clicks,
		})
	}

	return buckets, nil
}
