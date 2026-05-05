package analytics

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
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
					ID:        "link-1",
					OwnerID:   "owner-1",
					Code:      "abc123",
					TargetURL: "https://example.com",
					Status:    links.StatusActive,
					CreatedAt: now,
					UpdatedAt: now,
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

func TestHandlerLinkAnalyticsUnauthorized(t *testing.T) {
	handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/links/"+testAnalyticsLinkID+"/analytics", nil)
	rr := httptest.NewRecorder()

	handler.LinkAnalytics(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestHandlerLinkAnalyticsSuccess(t *testing.T) {
	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)
	repo := fakeDashboardRepository{
		getLinkByIDAndOwnerFn: func(_ context.Context, linkID, ownerID string) (links.Link, error) {
			if linkID != testAnalyticsLinkID {
				t.Fatalf("expected %s, got %s", testAnalyticsLinkID, linkID)
			}
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return links.Link{ID: testAnalyticsLinkID, OwnerID: "owner-1"}, nil
		},
		countLinkClicksFn: func(_ context.Context, _ string, _, _ time.Time) (int64, error) {
			return 4, nil
		},
		countLinkClicksSinceFn: func(_ context.Context, _ string, _ time.Time) (int64, error) {
			return 2, nil
		},
		lastLinkClickedAtFn: func(_ context.Context, _ string, _, _ time.Time) (*time.Time, error) {
			return &now, nil
		},
		listLinkClickSeriesFn: func(_ context.Context, _ string, _, _ time.Time, groupBy string) ([]TimeSeriesBucket, error) {
			if groupBy != GroupByHour {
				t.Fatalf("expected group by hour, got %s", groupBy)
			}
			return []TimeSeriesBucket{
				{PeriodStart: now.Truncate(time.Hour), Clicks: 4},
			}, nil
		},
	}
	service := NewService(repo, "https://tracklink.example.com")
	service.now = func() time.Time { return now }
	handler := NewHandler(service)
	req := newLinkAnalyticsRequest("/api/v1/links/"+testAnalyticsLinkID+"/analytics?groupBy=hour", testAnalyticsLinkID)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.LinkAnalytics(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp LinkAnalyticsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.LinkID != testAnalyticsLinkID {
		t.Fatalf("expected %s, got %s", testAnalyticsLinkID, resp.LinkID)
	}
	if resp.TotalClicks != 4 {
		t.Fatalf("expected total clicks 4, got %d", resp.TotalClicks)
	}
	if resp.ClicksLast24 != 2 {
		t.Fatalf("expected clicksLast24h 2, got %d", resp.ClicksLast24)
	}
	if resp.LastClickedAt == nil || *resp.LastClickedAt != now.Format(time.RFC3339) {
		t.Fatalf("unexpected lastClickedAt: %v", resp.LastClickedAt)
	}
	if len(resp.Series) != 1 || resp.Series[0].Clicks != 4 {
		t.Fatalf("unexpected series: %+v", resp.Series)
	}
}

func TestHandlerLinkAnalyticsPeriodFilter(t *testing.T) {
	from := time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC)
	to := time.Date(2026, 5, 2, 10, 30, 0, 0, time.UTC)
	var capturedFrom time.Time
	var capturedTo time.Time
	repo := fakeDashboardRepository{
		countLinkClicksFn: func(_ context.Context, _ string, fromArg, toArg time.Time) (int64, error) {
			capturedFrom = fromArg
			capturedTo = toArg
			return 1, nil
		},
	}
	handler := NewHandler(NewService(repo, "https://tracklink.example.com"))
	req := newLinkAnalyticsRequest("/api/v1/links/"+testAnalyticsLinkID+"/analytics?from=2026-04-28T09:00:00Z&to=2026-05-02T10:30:00Z&groupBy=day", testAnalyticsLinkID)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.LinkAnalytics(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if !capturedFrom.Equal(from) {
		t.Fatalf("expected from %s, got %s", from, capturedFrom)
	}
	if !capturedTo.Equal(to) {
		t.Fatalf("expected to %s, got %s", to, capturedTo)
	}
}

func TestHandlerLinkAnalyticsInvalidGroupBy(t *testing.T) {
	handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
	req := newLinkAnalyticsRequest("/api/v1/links/"+testAnalyticsLinkID+"/analytics?groupBy=week", testAnalyticsLinkID)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.LinkAnalytics(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandlerLinkAnalyticsInvalidPeriod(t *testing.T) {
	tests := []struct {
		name   string
		target string
	}{
		{
			name:   "invalid from",
			target: "/api/v1/links/" + testAnalyticsLinkID + "/analytics?from=not-a-date",
		},
		{
			name:   "invalid to",
			target: "/api/v1/links/" + testAnalyticsLinkID + "/analytics?to=not-a-date",
		},
		{
			name:   "from after to",
			target: "/api/v1/links/" + testAnalyticsLinkID + "/analytics?from=2026-05-03T00:00:00Z&to=2026-05-02T00:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
			req := newLinkAnalyticsRequest(tt.target, testAnalyticsLinkID)
			req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
				UserID: "owner-1",
				Role:   "customer",
			}))
			rr := httptest.NewRecorder()

			handler.LinkAnalytics(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
			}
		})
	}
}

func TestHandlerLinkAnalyticsNotFound(t *testing.T) {
	repo := fakeDashboardRepository{
		getLinkByIDAndOwnerFn: func(_ context.Context, _, _ string) (links.Link, error) {
			return links.Link{}, ErrLinkNotFound
		},
	}
	handler := NewHandler(NewService(repo, "https://tracklink.example.com"))
	req := newLinkAnalyticsRequest("/api/v1/links/"+testAnalyticsMissingLinkID+"/analytics", testAnalyticsMissingLinkID)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.LinkAnalytics(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandlerLinkAnalyticsInvalidLinkID(t *testing.T) {
	handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
	req := newLinkAnalyticsRequest("/api/v1/links/not-a-uuid/analytics", "not-a-uuid")
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.LinkAnalytics(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := resp.Error.Fields["linkId"]; !ok {
		t.Fatalf("expected linkId field error, got %v", resp.Error.Fields)
	}
}

func TestHandlerRecentClicksUnauthorized(t *testing.T) {
	handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/links/"+testAnalyticsLinkID+"/clicks", nil)
	rr := httptest.NewRecorder()

	handler.RecentClicks(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestHandlerRecentClicksSuccess(t *testing.T) {
	clickedAt := time.Date(2026, 5, 2, 13, 0, 0, 0, time.UTC)
	referrer := "https://t.me/example"
	userAgent := "Mozilla/5.0"
	repo := fakeDashboardRepository{
		listRecentClicksFn: func(_ context.Context, linkID string, limit int) ([]ClickEvent, error) {
			if linkID != testAnalyticsLinkID {
				t.Fatalf("expected %s, got %s", testAnalyticsLinkID, linkID)
			}
			if limit != 2 {
				t.Fatalf("expected limit 2, got %d", limit)
			}
			return []ClickEvent{
				{
					ID:        "click-1",
					LinkID:    testAnalyticsLinkID,
					ClickedAt: clickedAt,
					Referrer:  &referrer,
					UserAgent: &userAgent,
				},
			}, nil
		},
	}
	handler := NewHandler(NewService(repo, "https://tracklink.example.com"))
	req := newRecentClicksRequest("/api/v1/links/"+testAnalyticsLinkID+"/clicks?limit=2", testAnalyticsLinkID)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.RecentClicks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp RecentClicksResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected one click event, got %d", len(resp.Items))
	}
	item := resp.Items[0]
	if item.ID != "click-1" || item.LinkID != testAnalyticsLinkID {
		t.Fatalf("unexpected click event identity: %+v", item)
	}
	if item.ClickedAt != clickedAt.Format(time.RFC3339) {
		t.Fatalf("expected clickedAt %s, got %s", clickedAt.Format(time.RFC3339), item.ClickedAt)
	}
	if item.Referrer == nil || *item.Referrer != referrer {
		t.Fatalf("unexpected referrer: %v", item.Referrer)
	}
	if item.UserAgent == nil || *item.UserAgent != userAgent {
		t.Fatalf("unexpected userAgent: %v", item.UserAgent)
	}
}

func TestHandlerRecentClicksInvalidLimit(t *testing.T) {
	tests := []struct {
		name   string
		target string
	}{
		{
			name:   "not integer",
			target: "/api/v1/links/" + testAnalyticsLinkID + "/clicks?limit=abc",
		},
		{
			name:   "zero",
			target: "/api/v1/links/" + testAnalyticsLinkID + "/clicks?limit=0",
		},
		{
			name:   "above max",
			target: "/api/v1/links/" + testAnalyticsLinkID + "/clicks?limit=101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
			req := newRecentClicksRequest(tt.target, testAnalyticsLinkID)
			req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
				UserID: "owner-1",
				Role:   "customer",
			}))
			rr := httptest.NewRecorder()

			handler.RecentClicks(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
			}
		})
	}
}

func TestHandlerRecentClicksNotFound(t *testing.T) {
	repo := fakeDashboardRepository{
		getLinkByIDAndOwnerFn: func(_ context.Context, _, _ string) (links.Link, error) {
			return links.Link{}, ErrLinkNotFound
		},
	}
	handler := NewHandler(NewService(repo, "https://tracklink.example.com"))
	req := newRecentClicksRequest("/api/v1/links/"+testAnalyticsMissingLinkID+"/clicks", testAnalyticsMissingLinkID)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.RecentClicks(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandlerRecentClicksInvalidLinkID(t *testing.T) {
	handler := NewHandler(NewService(fakeDashboardRepository{}, "https://tracklink.example.com"))
	req := newRecentClicksRequest("/api/v1/links/not-a-uuid/clicks", "not-a-uuid")
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.RecentClicks(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := resp.Error.Fields["linkId"]; !ok {
		t.Fatalf("expected linkId field error, got %v", resp.Error.Fields)
	}
}

func TestHandlerRecentClicksInternalError(t *testing.T) {
	repo := fakeDashboardRepository{
		listRecentClicksFn: func(_ context.Context, _ string, _ int) ([]ClickEvent, error) {
			return nil, errors.New("db down")
		},
	}
	handler := NewHandler(NewService(repo, "https://tracklink.example.com"))
	req := newRecentClicksRequest("/api/v1/links/"+testAnalyticsLinkID+"/clicks", testAnalyticsLinkID)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-1", session.SessionData{
		UserID: "owner-1",
		Role:   "customer",
	}))
	rr := httptest.NewRecorder()

	handler.RecentClicks(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func newLinkAnalyticsRequest(target, linkID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("linkId", linkID)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
}

func newRecentClicksRequest(target, linkID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("linkId", linkID)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
}
