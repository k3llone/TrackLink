package redirect

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	findByCodeOrAliasFn func(ctx context.Context, code string) (Link, error)
}

func (f fakeRepository) FindByCodeOrAlias(ctx context.Context, code string) (Link, error) {
	if f.findByCodeOrAliasFn == nil {
		return Link{}, ErrLinkNotFound
	}

	return f.findByCodeOrAliasFn(ctx, code)
}

func TestServiceResolveByCode(t *testing.T) {
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != "abc123" {
				t.Fatalf("expected code abc123, got %s", code)
			}
			return Link{
				ID:        "link-1",
				Code:      "abc123",
				TargetURL: "https://example.com/landing",
				Status:    "active",
			}, nil
		},
	}

	service := NewService(repo)
	link, err := service.Resolve(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if link.TargetURL != "https://example.com/landing" {
		t.Fatalf("expected target url to match, got %s", link.TargetURL)
	}
}

func TestServiceResolveByCustomAlias(t *testing.T) {
	alias := "promo"
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != alias {
				t.Fatalf("expected alias %s, got %s", alias, code)
			}
			return Link{
				ID:          "link-2",
				Code:        "a1b2c3",
				CustomAlias: &alias,
				TargetURL:   "https://example.com/promo",
				Status:      "active",
			}, nil
		},
	}

	service := NewService(repo)
	link, err := service.Resolve(context.Background(), "promo")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if link.CustomAlias == nil || *link.CustomAlias != alias {
		t.Fatalf("expected custom alias %s, got %v", alias, link.CustomAlias)
	}
}

func TestServiceResolveNotFound(t *testing.T) {
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, _ string) (Link, error) {
			return Link{}, ErrLinkNotFound
		},
	}

	service := NewService(repo)
	_, err := service.Resolve(context.Background(), "missing")
	if !errors.Is(err, ErrLinkNotFound) {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}
}
