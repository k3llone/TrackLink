package redirect

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
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
