package redirect

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrLinkNotFound = errors.New("link not found")

type Link struct {
	ID          string
	Code        string
	CustomAlias *string
	TargetURL   string
	Status      string
	DeletedAt   *time.Time
}

type Repository interface {
	FindByCodeOrAlias(ctx context.Context, code string) (Link, error)
	TouchActiveLink(ctx context.Context, linkID string, clickedAt time.Time) error
}

type GormRepository struct {
	db *gorm.DB
}

type linkModel struct {
	ID          string
	Code        string
	CustomAlias *string
	TargetURL   string
	Status      string
	DeletedAt   *time.Time
}

func (linkModel) TableName() string {
	return "links"
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindByCodeOrAlias(ctx context.Context, code string) (Link, error) {
	if r.db == nil {
		return Link{}, fmt.Errorf("find link by code or alias: db is nil")
	}

	var row linkModel
	err := r.db.WithContext(ctx).
		Where("code = ? OR custom_alias = ?", code, code).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, fmt.Errorf("find link by code or alias: %w", err)
	}

	return Link{
		ID:          row.ID,
		Code:        row.Code,
		CustomAlias: row.CustomAlias,
		TargetURL:   row.TargetURL,
		Status:      row.Status,
		DeletedAt:   row.DeletedAt,
	}, nil
}

func (r *GormRepository) TouchActiveLink(ctx context.Context, linkID string, clickedAt time.Time) error {
	if r.db == nil {
		return fmt.Errorf("touch active link: db is nil")
	}

	result := r.db.WithContext(ctx).
		Model(&linkModel{}).
		Where("id = ? AND status = ? AND deleted_at IS NULL", linkID, StatusActive).
		Clauses(clause.Returning{}).
		Updates(map[string]any{
			"total_clicks":    gorm.Expr("total_clicks + 1"),
			"last_clicked_at": clickedAt.UTC(),
			"updated_at":      time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("touch active link: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrLinkNotFound
	}

	return nil
}
