package redirect

import (
	"context"
	"errors"
	"strings"
	"time"

	"tracklink/internal/modules/analytics"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusBlocked  = "blocked"
	StatusDeleted  = "deleted"
)

type ResultKind string

const (
	ResultKindRedirect    ResultKind = "redirect"
	ResultKindNotFound    ResultKind = "not_found"
	ResultKindUnavailable ResultKind = "unavailable"
)

type RequestMeta struct {
	Referrer  string
	UserAgent string
	ClickedAt time.Time
}

type ResolveResult struct {
	Kind      ResultKind
	TargetURL string
	Status    string
	Deleted   bool
}

type Service struct {
	repo          Repository
	analyticsRepo analytics.Repository
}

func NewService(repo Repository, analyticsRepo analytics.Repository) *Service {
	return &Service{
		repo:          repo,
		analyticsRepo: analyticsRepo,
	}
}

func (s *Service) ResolveAndTrack(ctx context.Context, code string, meta RequestMeta) (ResolveResult, error) {
	normalizedCode := strings.TrimSpace(code)
	if normalizedCode == "" {
		return ResolveResult{Kind: ResultKindNotFound}, nil
	}

	link, err := s.repo.FindByCodeOrAlias(ctx, normalizedCode)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return ResolveResult{Kind: ResultKindNotFound}, nil
		}
		return ResolveResult{}, err
	}

	if link.Status != StatusActive || link.DeletedAt != nil {
		return unavailableResult(link), nil
	}

	clickedAt := meta.ClickedAt.UTC()
	if clickedAt.IsZero() {
		clickedAt = time.Now().UTC()
	}

	if err := s.repo.TouchActiveLink(ctx, link.ID, clickedAt); err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			return unavailableResult(link), nil
		}
		return ResolveResult{}, err
	}

	if s.analyticsRepo != nil {
		if err := s.analyticsRepo.CreateClickEvent(ctx, analytics.CreateClickEventParams{
			LinkID:    link.ID,
			ClickedAt: clickedAt,
			Referrer:  normalizeNullableString(meta.Referrer),
			UserAgent: normalizeNullableString(meta.UserAgent),
		}); err != nil {
			// Analytics failure must not block redirect flow.
		}
	}

	return ResolveResult{
		Kind:      ResultKindRedirect,
		TargetURL: link.TargetURL,
		Status:    link.Status,
	}, nil
}

func unavailableResult(link Link) ResolveResult {
	return ResolveResult{
		Kind:    ResultKindUnavailable,
		Status:  link.Status,
		Deleted: link.Status == StatusDeleted || link.DeletedAt != nil,
	}
}

func normalizeNullableString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
