package analytics

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tracklink/internal/modules/links"
	"tracklink/internal/platform/session"
	"tracklink/internal/shared"
)

func TestHandlerDashboardUnauthorized(t *testing.T) {
	handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	rr := httptest.NewRecorder()

	handler.Dashboard(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestHandlerDashboardSuccess(t *testing.T) {
	repo := fakeDashboardRepository{
		countTotalLinksFn: func(_ context.Context, ownerID string) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return 4, nil
		},
		countActiveLinksFn: func(_ context.Context, ownerID string) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return 3, nil
		},
		sumTotalClicksFn: func(_ context.Context, ownerID string) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return 7, nil
		},
		countClicksSinceFn: func(_ context.Context, ownerID string, since time.Time) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			if since.IsZero() {
				t.Fatal("expected non-zero since")
			}
			return 2, nil
		},
		listRecentLinksFn: func(_ context.Context, _ string, _ int) ([]links.Link, error) {
			now := time.Date(2026, 5, 2, 9, 0, 0, 0, time.UTC)
			return []links.Link{
				{
					ID:         "link-1",
					OwnerID:    "owner-1",
					Code:       "abc123",
					TargetURL:  "https://example.com",
					Status:     links.StatusActive,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			}, nil
		},
	}
	handler := NewHandler(NewService(repo, "https://tracklink.example.com"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.Dashboard(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp DashboardResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.TotalClicks != 7 {
		t.Fatalf("expected total clicks 7, got %d", resp.TotalClicks)
	}
	if resp.TotalLinks != 4 {
		t.Fatalf("expected total links 4, got %d", resp.TotalLinks)
	}
	if resp.ActiveLinks != 3 {
		t.Fatalf("expected active links 3, got %d", resp.ActiveLinks)
	}
	if resp.ClicksLast24 != 2 {
		t.Fatalf("expected clicksLast24h 2, got %d", resp.ClicksLast24)
	}
	if len(resp.RecentLinks) != 1 {
		t.Fatalf("expected one recent link, got %d", len(resp.RecentLinks))
	}
}

func TestHandlerDashboardInternalError(t *testing.T) {
	repo := fakeDashboardRepository{
		sumTotalClicksFn: func(_ context.Context, _ string) (int64, error) {
			return 0, errors.New("db down")
		},
	}
	handler := NewHandler(NewService(repo, "https://tracklink.example.com"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.Dashboard(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
