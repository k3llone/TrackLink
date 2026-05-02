package redirect

import (
	"context"
	"testing"
	"time"

	"tracklink/internal/modules/analytics"
)

type fakeRepository struct {
	findByCodeOrAliasFn func(ctx context.Context, code string) (Link, error)
	touchActiveLinkFn   func(ctx context.Context, linkID string, clickedAt time.Time) error
}

func (f fakeRepository) FindByCodeOrAlias(ctx context.Context, code string) (Link, error) {
	if f.findByCodeOrAliasFn == nil {
		return Link{}, ErrLinkNotFound
	}

	return f.findByCodeOrAliasFn(ctx, code)
}

type fakeAnalyticsRepository struct {
	createClickEventFn func(ctx context.Context, event analytics.CreateClickEventParams) error
}

func (f fakeAnalyticsRepository) CreateClickEvent(ctx context.Context, event analytics.CreateClickEventParams) error {
	if f.createClickEventFn == nil {
		return nil
	}

	return f.createClickEventFn(ctx, event)
}

func (f fakeRepository) TouchActiveLink(ctx context.Context, linkID string, clickedAt time.Time) error {
	if f.touchActiveLinkFn == nil {
		return nil
	}

	return f.touchActiveLinkFn(ctx, linkID, clickedAt)
}

func TestServiceResolveByCodeTracksAndRedirectsActive(t *testing.T) {
	clickedAt := time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC)
	steps := make([]string, 0, 2)
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, event analytics.CreateClickEventParams) error {
			steps = append(steps, "event")
			if event.LinkID != "link-1" {
				t.Fatalf("expected link id link-1, got %s", event.LinkID)
			}
			if !event.ClickedAt.Equal(clickedAt) {
				t.Fatalf("expected clickedAt %v, got %v", clickedAt, event.ClickedAt)
			}
			return nil
		},
	}
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %s", code)
			}
			return Link{
				ID:        "link-1",
				Code:      "abc123",
				TargetURL: "https://example.com/landing",
				Status:    StatusActive,
			}, nil
		},
		touchActiveLinkFn: func(_ context.Context, linkID string, gotClickedAt time.Time) error {
			steps = append(steps, "touch")
			if linkID != "link-1" {
				t.Fatalf("expected link id link-1, got %s", linkID)
			}
			if !gotClickedAt.Equal(clickedAt) {
				t.Fatalf("expected clickedAt %v, got %v", clickedAt, gotClickedAt)
			}
			return nil
		},
	}

	service := NewService(repo, analyticsRepo)
	result, err := service.ResolveAndTrack(context.Background(), "abc123", RequestMeta{
		Referrer:  "https://source.example.com",
		UserAgent: "test-agent",
		ClickedAt: clickedAt,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Kind != ResultKindRedirect {
		t.Fatalf("expected redirect result, got %s", result.Kind)
	}
	if result.TargetURL != "https://example.com/landing" {
		t.Fatalf("expected target url to match, got %s", result.TargetURL)
	}
	if len(steps) != 2 || steps[0] != "event" || steps[1] != "touch" {
		t.Fatalf("expected click event before touch, got %v", steps)
	}
}

func TestServiceResolveByCustomAliasAsNotFoundWhenMissing(t *testing.T) {
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != "promo" {
				t.Fatalf("expected alias promo, got %s", code)
			}
			return Link{}, ErrLinkNotFound
		},
	}
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, _ analytics.CreateClickEventParams) error {
			t.Fatal("create click event should not be called for not found")
			return nil
		},
	}

	service := NewService(repo, analyticsRepo)
	result, err := service.ResolveAndTrack(context.Background(), "promo", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Kind != ResultKindNotFound {
		t.Fatalf("expected not found, got %s", result.Kind)
	}
}

func TestServiceResolveTracksUnavailableWithoutTouch(t *testing.T) {
	clicked := false
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, event analytics.CreateClickEventParams) error {
			clicked = true
			if event.LinkID != "link-3" {
				t.Fatalf("expected link id link-3, got %s", event.LinkID)
			}
			return nil
		},
	}
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != "blocked-link" {
				t.Fatalf("expected code blocked-link, got %s", code)
			}
			return Link{
				ID:        "link-3",
				Code:      "blocked-link",
				TargetURL: "https://example.com/blocked",
				Status:    StatusBlocked,
			}, nil
		},
		touchActiveLinkFn: func(_ context.Context, _ string, _ time.Time) error {
			t.Fatal("touch active link must not be called for blocked link")
			return nil
		},
	}

	service := NewService(repo, analyticsRepo)
	result, err := service.ResolveAndTrack(context.Background(), "blocked-link", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Kind != ResultKindUnavailable {
		t.Fatalf("expected unavailable result, got %s", result.Kind)
	}
	if result.Status != StatusBlocked {
		t.Fatalf("expected blocked status, got %s", result.Status)
	}
	if !clicked {
		t.Fatal("expected click event to be created for unavailable link")
	}
}

func TestServiceResolveKeepsRedirectWhenAnalyticsSaveFails(t *testing.T) {
	touched := false
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, _ string) (Link, error) {
			return Link{
				ID:        "link-10",
				Code:      "go",
				TargetURL: "https://example.com/go",
				Status:    StatusActive,
			}, nil
		},
		touchActiveLinkFn: func(_ context.Context, linkID string, _ time.Time) error {
			touched = true
			if linkID != "link-10" {
				t.Fatalf("expected link id link-10, got %s", linkID)
			}
			return nil
		},
	}
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, _ analytics.CreateClickEventParams) error {
			return context.DeadlineExceeded
		},
	}

	service := NewService(repo, analyticsRepo)
	result, err := service.ResolveAndTrack(context.Background(), "go", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Kind != ResultKindRedirect {
		t.Fatalf("expected redirect result, got %s", result.Kind)
	}
	if result.TargetURL != "https://example.com/go" {
		t.Fatalf("expected target url to match, got %s", result.TargetURL)
	}
	if !touched {
		t.Fatal("expected active link click stats to be updated")
	}
}
