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
