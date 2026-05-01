package links

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, link *Link) error
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByCustomAlias(ctx context.Context, customAlias string) (bool, error)
	ListByOwner(ctx context.Context, filter ListLinksFilter) ([]Link, int64, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, link *Link) error {
	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		if isUniqueConstraintViolation(err, "idx_links_custom_alias_unique") {
			return ErrAliasAlreadyExists
		}
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

func (r *GormRepository) ExistsByCustomAlias(ctx context.Context, customAlias string) (bool, error) {
	var link Link
	err := r.db.WithContext(ctx).
		Select("id").
		Where("custom_alias = ?", customAlias).
		First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("check custom alias existence: %w", err)
	}

	return true, nil
}

func isUniqueConstraintViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraintName {
		return true
	}

	return false
}

func (r *GormRepository) ListByOwner(ctx context.Context, filter ListLinksFilter) ([]Link, int64, error) {
	base := r.db.WithContext(ctx).Model(&Link{}).Where("owner_id = ?", filter.OwnerID)

	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
		if filter.Status != StatusDeleted {
			base = base.Where("deleted_at IS NULL")
		}
	} else {
		base = base.Where("status <> ?", StatusDeleted).Where("deleted_at IS NULL")
	}

	if filter.Q != "" {
		like := "%" + filter.Q + "%"
		base = base.Where("(code ILIKE ? OR custom_alias ILIKE ? OR target_url ILIKE ?)", like, like, like)
	}

	var totalItems int64
	if err := base.Count(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("count links by owner: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	var links []Link
	if err := base.
		Order("created_at DESC").
		Limit(filter.PageSize).
		Offset(offset).
		Find(&links).Error; err != nil {
		return nil, 0, fmt.Errorf("list links by owner: %w", err)
	}

	return links, totalItems, nil
}
