package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tracklink/internal/modules/links"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, filter ListLinksFilter) ([]links.Link, int64, error)
	GetByID(ctx context.Context, linkID string) (links.Link, error)
	UpdateStatus(ctx context.Context, linkID, status string) (links.Link, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context, filter ListLinksFilter) ([]links.Link, int64, error) {
	base := r.db.WithContext(ctx).
		Model(&links.Link{}).
		Where("status <> ?", links.StatusDeleted).
		Where("deleted_at IS NULL")

	if filter.Q != "" {
		like := "%" + filter.Q + "%"
		base = base.Where("(id::text ILIKE ? OR code ILIKE ? OR custom_alias ILIKE ?)", like, like, like)
	}

	var totalItems int64
	if err := base.Count(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("count admin links: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	var items []links.Link
	if err := base.
		Order("created_at DESC").
		Limit(filter.PageSize).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list admin links: %w", err)
	}

	return items, totalItems, nil
}

func (r *GormRepository) GetByID(ctx context.Context, linkID string) (links.Link, error) {
	var link links.Link
	err := r.db.WithContext(ctx).
		Where("id = ?", linkID).
		Where("status <> ?", links.StatusDeleted).
		Where("deleted_at IS NULL").
		First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return links.Link{}, ErrLinkNotFound
		}
		return links.Link{}, fmt.Errorf("get admin link by id: %w", err)
	}

	return link, nil
}

func (r *GormRepository) UpdateStatus(ctx context.Context, linkID, status string) (links.Link, error) {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&links.Link{}).
		Where("id = ?", linkID).
		Where("status <> ?", links.StatusDeleted).
		Where("deleted_at IS NULL").
		Updates(map[string]any{
			"status":     status,
			"updated_at": now,
		})
	if result.Error != nil {
		return links.Link{}, fmt.Errorf("update admin link status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return links.Link{}, ErrLinkNotFound
	}

	return r.GetByID(ctx, linkID)
}
