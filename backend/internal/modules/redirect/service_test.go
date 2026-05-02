package redirect

import (
	"context"
	"testing"
	"time"
)

type fakeRepository struct {
	findByCodeOrAliasFn func(ctx context.Context, code string) (Link, error)
	createClickEventFn  func(ctx context.Context, event ClickEvent) error
	touchActiveLinkFn   func(ctx context.Context, linkID string, clickedAt time.Time) error
}

func (f fakeRepository) FindByCodeOrAlias(ctx context.Context, code string) (Link, error) {
	if f.findByCodeOrAliasFn == nil {
		return Link{}, ErrLinkNotFound
	}

	return f.findByCodeOrAliasFn(ctx, code)
}

func (f fakeRepository) CreateClickEvent(ctx context.Context, event ClickEvent) error {
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
		createClickEventFn: func(_ context.Context, event ClickEvent) error {
			steps = append(steps, "event")
			if event.LinkID != "link-1" {
				t.Fatalf("expected link id link-1, got %s", event.LinkID)
			}
			if !event.ClickedAt.Equal(clickedAt) {
				t.Fatalf("expected clickedAt %v, got %v", clickedAt, event.ClickedAt)
			}
			return nil
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

	service := NewService(repo)
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
		createClickEventFn: func(_ context.Context, _ ClickEvent) error {
			t.Fatal("create click event should not be called for not found")
			return nil
		},
	}

	service := NewService(repo)
	result, err := service.ResolveAndTrack(context.Background(), "promo", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Kind != ResultKindNotFound {
		t.Fatalf("expected not found, got %s", result.Kind)
	}
}
