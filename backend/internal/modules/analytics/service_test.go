package analytics

import (
	"context"
	"testing"
	"time"

	"tracklink/internal/modules/links"
)

type fakeDashboardRepository struct {
	sumTotalClicksFn func(ctx context.Context, ownerID string) (int64, error)
	listRecentLinksFn func(ctx context.Context, ownerID string, limit int) ([]links.Link, error)
}

func (f fakeDashboardRepository) SumTotalClicks(ctx context.Context, ownerID string) (int64, error) {
	if f.sumTotalClicksFn == nil {
		return 0, nil
	}
	return f.sumTotalClicksFn(ctx, ownerID)
}

func (f fakeDashboardRepository) ListRecentLinks(ctx context.Context, ownerID string, limit int) ([]links.Link, error) {
	if f.listRecentLinksFn == nil {
		return nil, nil
	}
	return f.listRecentLinksFn(ctx, ownerID, limit)
}

func TestServiceLoadDashboardReturnsTotalClicksAndRecentLinks(t *testing.T) {
	createdAt := time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(2 * time.Hour)
	lastClickedAt := createdAt.Add(6 * time.Hour)
	repo := fakeDashboardRepository{
		sumTotalClicksFn: func(_ context.Context, ownerID string) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return 42, nil
		},
		listRecentLinksFn: func(_ context.Context, ownerID string, limit int) ([]links.Link, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			if limit != 5 {
				t.Fatalf("expected limit 5, got %d", limit)
			}
			return []links.Link{
				{
					ID:            "link-1",
					OwnerID:       "owner-1",
					Code:          "abc123",
					TargetURL:     "https://example.com",
					Status:        links.StatusActive,
					TotalClicks:   42,
					LastClickedAt: &lastClickedAt,
					CreatedAt:     createdAt,
					UpdatedAt:     updatedAt,
				},
			}, nil
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	resp, fields, err := service.LoadDashboard(context.Background(), "owner-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(fields) != 0 {
		t.Fatalf("expected empty fields, got %v", fields)
	}
	if resp.TotalClicks != 42 {
		t.Fatalf("expected total clicks 42, got %d", resp.TotalClicks)
	}
	if len(resp.RecentLinks) != 1 {
		t.Fatalf("expected one recent link, got %d", len(resp.RecentLinks))
	}
	if resp.RecentLinks[0].ShortURL != "https://tracklink.example.com/abc123" {
		t.Fatalf("unexpected short url: %s", resp.RecentLinks[0].ShortURL)
	}
	if resp.RecentLinks[0].TotalClicks != 42 {
		t.Fatalf("expected link total clicks 42, got %d", resp.RecentLinks[0].TotalClicks)
	}
	if resp.TotalLinks != 0 || resp.ActiveLinks != 0 || resp.ClicksLast24 != 0 {
		t.Fatalf("expected other aggregates to be zero, got %+v", resp)
	}
}

func TestServiceLoadDashboardValidation(t *testing.T) {
	service := NewService(fakeDashboardRepository{}, "https://tracklink.example.com")

	_, fields, err := service.LoadDashboard(context.Background(), "")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if fields["ownerId"] == "" {
		t.Fatalf("expected ownerId validation field, got %v", fields)
	}
}
