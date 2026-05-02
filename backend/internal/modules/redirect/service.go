package redirect

import (
	"context"
	"errors"
	"strings"
	"time"
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
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
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

	clickedAt := meta.ClickedAt.UTC()
	if clickedAt.IsZero() {
		clickedAt = time.Now().UTC()
	}
	if err := s.repo.CreateClickEvent(ctx, ClickEvent{
		LinkID:    link.ID,
		ClickedAt: clickedAt,
		Referrer:  strings.TrimSpace(meta.Referrer),
		UserAgent: strings.TrimSpace(meta.UserAgent),
	}); err != nil {
		return ResolveResult{}, err
	}

	if link.Status == StatusActive && link.DeletedAt == nil {
		if err := s.repo.TouchActiveLink(ctx, link.ID, clickedAt); err != nil {
			return ResolveResult{}, err
		}
		return ResolveResult{
			Kind:      ResultKindRedirect,
			TargetURL: link.TargetURL,
			Status:    link.Status,
		}, nil
	}

	return ResolveResult{
		Kind:    ResultKindUnavailable,
		Status:  link.Status,
		Deleted: link.Status == StatusDeleted || link.DeletedAt != nil,
	}, nil
}
