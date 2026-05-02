package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"tracklink/internal/modules/links"
)

type fakeDashboardRepository struct {
	countTotalLinksFn      func(ctx context.Context, ownerID string) (int64, error)
	countActiveLinksFn     func(ctx context.Context, ownerID string) (int64, error)
	sumTotalClicksFn       func(ctx context.Context, ownerID string) (int64, error)
	countClicksSinceFn     func(ctx context.Context, ownerID string, since time.Time) (int64, error)
	listRecentLinksFn      func(ctx context.Context, ownerID string, limit int) ([]links.Link, error)
	getLinkByIDAndOwnerFn  func(ctx context.Context, linkID, ownerID string) (links.Link, error)
	countLinkClicksFn      func(ctx context.Context, linkID string, from, to time.Time) (int64, error)
	countLinkClicksSinceFn func(ctx context.Context, linkID string, since time.Time) (int64, error)
	lastLinkClickedAtFn    func(ctx context.Context, linkID string, from, to time.Time) (*time.Time, error)
	listLinkClickSeriesFn  func(ctx context.Context, linkID string, from, to time.Time, groupBy string) ([]TimeSeriesBucket, error)
	listRecentClicksFn     func(ctx context.Context, linkID string, limit int) ([]ClickEvent, error)
}

func (f fakeDashboardRepository) CountTotalLinks(ctx context.Context, ownerID string) (int64, error) {
	if f.countTotalLinksFn == nil {
		return 0, nil
	}
	return f.countTotalLinksFn(ctx, ownerID)
}

func (f fakeDashboardRepository) CountActiveLinks(ctx context.Context, ownerID string) (int64, error) {
	if f.countActiveLinksFn == nil {
		return 0, nil
	}
	return f.countActiveLinksFn(ctx, ownerID)
}

func (f fakeDashboardRepository) SumTotalClicks(ctx context.Context, ownerID string) (int64, error) {
	if f.sumTotalClicksFn == nil {
		return 0, nil
	}
	return f.sumTotalClicksFn(ctx, ownerID)
}

func (f fakeDashboardRepository) CountClicksSince(ctx context.Context, ownerID string, since time.Time) (int64, error) {
	if f.countClicksSinceFn == nil {
		return 0, nil
	}
	return f.countClicksSinceFn(ctx, ownerID, since)
}

func (f fakeDashboardRepository) ListRecentLinks(ctx context.Context, ownerID string, limit int) ([]links.Link, error) {
	if f.listRecentLinksFn == nil {
		return nil, nil
	}
	return f.listRecentLinksFn(ctx, ownerID, limit)
}

func (f fakeDashboardRepository) GetLinkByIDAndOwner(ctx context.Context, linkID, ownerID string) (links.Link, error) {
	if f.getLinkByIDAndOwnerFn == nil {
		return links.Link{ID: linkID, OwnerID: ownerID}, nil
	}
	return f.getLinkByIDAndOwnerFn(ctx, linkID, ownerID)
}

func (f fakeDashboardRepository) CountLinkClicks(ctx context.Context, linkID string, from, to time.Time) (int64, error) {
	if f.countLinkClicksFn == nil {
		return 0, nil
	}
	return f.countLinkClicksFn(ctx, linkID, from, to)
}

func (f fakeDashboardRepository) CountLinkClicksSince(ctx context.Context, linkID string, since time.Time) (int64, error) {
	if f.countLinkClicksSinceFn == nil {
		return 0, nil
	}
	return f.countLinkClicksSinceFn(ctx, linkID, since)
}

func (f fakeDashboardRepository) LastLinkClickedAt(ctx context.Context, linkID string, from, to time.Time) (*time.Time, error) {
	if f.lastLinkClickedAtFn == nil {
		return nil, nil
	}
	return f.lastLinkClickedAtFn(ctx, linkID, from, to)
}

func (f fakeDashboardRepository) ListLinkClickSeries(ctx context.Context, linkID string, from, to time.Time, groupBy string) ([]TimeSeriesBucket, error) {
	if f.listLinkClickSeriesFn == nil {
		return nil, nil
	}
	return f.listLinkClickSeriesFn(ctx, linkID, from, to, groupBy)
}

func (f fakeDashboardRepository) ListRecentClicks(ctx context.Context, linkID string, limit int) ([]ClickEvent, error) {
	if f.listRecentClicksFn == nil {
		return nil, nil
	}
	return f.listRecentClicksFn(ctx, linkID, limit)
}

func TestServiceLoadDashboardReturnsTotalClicksAndRecentLinks(t *testing.T) {
	createdAt := time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(2 * time.Hour)
	lastClickedAt := createdAt.Add(6 * time.Hour)
	repo := fakeDashboardRepository{
		countTotalLinksFn: func(_ context.Context, ownerID string) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return 3, nil
		},
		countActiveLinksFn: func(_ context.Context, ownerID string) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return 2, nil
		},
		sumTotalClicksFn: func(_ context.Context, ownerID string) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return 42, nil
		},
		countClicksSinceFn: func(_ context.Context, ownerID string, since time.Time) (int64, error) {
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			if since.IsZero() {
				t.Fatal("expected non-zero since")
			}
			return 5, nil
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
	if resp.TotalLinks != 3 {
		t.Fatalf("expected total links 3, got %d", resp.TotalLinks)
	}
	if resp.ActiveLinks != 2 {
		t.Fatalf("expected active links 2, got %d", resp.ActiveLinks)
	}
	if resp.ClicksLast24 != 5 {
		t.Fatalf("expected clicksLast24h 5, got %d", resp.ClicksLast24)
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

func TestServiceLoadDashboardReturnsErrorWhenCountClicksSinceFails(t *testing.T) {
	repo := fakeDashboardRepository{
		sumTotalClicksFn: func(_ context.Context, _ string) (int64, error) {
			return 42, nil
		},
		countClicksSinceFn: func(_ context.Context, _ string, _ time.Time) (int64, error) {
			return 0, errors.New("count failed")
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	_, _, err := service.LoadDashboard(context.Background(), "owner-1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServiceLoadDashboardReturnsErrorWhenTotalLinksFails(t *testing.T) {
	repo := fakeDashboardRepository{
		countTotalLinksFn: func(_ context.Context, _ string) (int64, error) {
			return 0, errors.New("count links failed")
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	_, _, err := service.LoadDashboard(context.Background(), "owner-1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServiceLoadLinkAnalyticsReturnsDefaultPeriodAndSeries(t *testing.T) {
	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)
	lastClickedAt := now.Add(-2 * time.Hour)
	var capturedFrom time.Time
	var capturedTo time.Time
	var capturedGroupBy string
	repo := fakeDashboardRepository{
		getLinkByIDAndOwnerFn: func(_ context.Context, linkID, ownerID string) (links.Link, error) {
			if linkID != "link-1" {
				t.Fatalf("expected link-1, got %s", linkID)
			}
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return links.Link{ID: "link-1", OwnerID: "owner-1"}, nil
		},
		countLinkClicksFn: func(_ context.Context, linkID string, from, to time.Time) (int64, error) {
			if linkID != "link-1" {
				t.Fatalf("expected link-1, got %s", linkID)
			}
			capturedFrom = from
			capturedTo = to
			return 8, nil
		},
		countLinkClicksSinceFn: func(_ context.Context, linkID string, since time.Time) (int64, error) {
			if linkID != "link-1" {
				t.Fatalf("expected link-1, got %s", linkID)
			}
			expected := now.Add(-24 * time.Hour)
			if !since.Equal(expected) {
				t.Fatalf("expected since %s, got %s", expected, since)
			}
			return 3, nil
		},
		lastLinkClickedAtFn: func(_ context.Context, _ string, _, _ time.Time) (*time.Time, error) {
			return &lastClickedAt, nil
		},
		listLinkClickSeriesFn: func(_ context.Context, linkID string, from, to time.Time, groupBy string) ([]TimeSeriesBucket, error) {
			if linkID != "link-1" {
				t.Fatalf("expected link-1, got %s", linkID)
			}
			if !from.Equal(capturedFrom) || !to.Equal(capturedTo) {
				t.Fatal("expected series period to match click count period")
			}
			capturedGroupBy = groupBy
			return []TimeSeriesBucket{
				{PeriodStart: lastClickedAt.Truncate(time.Hour), Clicks: 5},
			}, nil
		},
	}
	service := NewService(repo, "https://tracklink.example.com")
	service.now = func() time.Time { return now }

	resp, fields, err := service.LoadLinkAnalytics(context.Background(), "owner-1", "link-1", LinkAnalyticsQuery{
		GroupBy: GroupByHour,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(fields) != 0 {
		t.Fatalf("expected empty fields, got %v", fields)
	}
	if !capturedFrom.Equal(now.Add(-defaultAnalyticsPeriod)) {
		t.Fatalf("expected default from %s, got %s", now.Add(-defaultAnalyticsPeriod), capturedFrom)
	}
	if !capturedTo.Equal(now) {
		t.Fatalf("expected default to %s, got %s", now, capturedTo)
	}
	if capturedGroupBy != GroupByHour {
		t.Fatalf("expected group by hour, got %s", capturedGroupBy)
	}
	if resp.LinkID != "link-1" {
		t.Fatalf("expected link ID link-1, got %s", resp.LinkID)
	}
	if resp.TotalClicks != 8 {
		t.Fatalf("expected total clicks 8, got %d", resp.TotalClicks)
	}
	if resp.ClicksLast24 != 3 {
		t.Fatalf("expected clicksLast24h 3, got %d", resp.ClicksLast24)
	}
	if resp.LastClickedAt == nil || *resp.LastClickedAt != lastClickedAt.Format(time.RFC3339) {
		t.Fatalf("unexpected lastClickedAt: %v", resp.LastClickedAt)
	}
	if len(resp.Series) != 1 {
		t.Fatalf("expected one series point, got %d", len(resp.Series))
	}
	if resp.Series[0].PeriodStart != lastClickedAt.Truncate(time.Hour).Format(time.RFC3339) || resp.Series[0].Clicks != 5 {
		t.Fatalf("unexpected series point: %+v", resp.Series[0])
	}
}

func TestServiceLoadLinkAnalyticsDefaultsGroupByDay(t *testing.T) {
	var capturedGroupBy string
	repo := fakeDashboardRepository{
		listLinkClickSeriesFn: func(_ context.Context, _ string, _, _ time.Time, groupBy string) ([]TimeSeriesBucket, error) {
			capturedGroupBy = groupBy
			return nil, nil
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	_, _, err := service.LoadLinkAnalytics(context.Background(), "owner-1", "link-1", LinkAnalyticsQuery{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if capturedGroupBy != GroupByDay {
		t.Fatalf("expected group by day, got %s", capturedGroupBy)
	}
}

func TestServiceLoadLinkAnalyticsReturnsNotFoundForForeignLink(t *testing.T) {
	repo := fakeDashboardRepository{
		getLinkByIDAndOwnerFn: func(_ context.Context, _, _ string) (links.Link, error) {
			return links.Link{}, ErrLinkNotFound
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	_, _, err := service.LoadLinkAnalytics(context.Background(), "owner-1", "link-1", LinkAnalyticsQuery{})
	if !errors.Is(err, ErrLinkNotFound) {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}
}

func TestServiceLoadLinkAnalyticsEmptyClicksResponse(t *testing.T) {
	repo := fakeDashboardRepository{
		countLinkClicksFn: func(_ context.Context, _ string, _, _ time.Time) (int64, error) {
			return 0, nil
		},
		countLinkClicksSinceFn: func(_ context.Context, _ string, _ time.Time) (int64, error) {
			return 0, nil
		},
		lastLinkClickedAtFn: func(_ context.Context, _ string, _, _ time.Time) (*time.Time, error) {
			return nil, nil
		},
		listLinkClickSeriesFn: func(_ context.Context, _ string, _, _ time.Time, _ string) ([]TimeSeriesBucket, error) {
			return []TimeSeriesBucket{}, nil
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	resp, _, err := service.LoadLinkAnalytics(context.Background(), "owner-1", "link-1", LinkAnalyticsQuery{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.TotalClicks != 0 || resp.ClicksLast24 != 0 {
		t.Fatalf("expected zero click counters, got %+v", resp)
	}
	if resp.LastClickedAt != nil {
		t.Fatalf("expected nil lastClickedAt, got %v", resp.LastClickedAt)
	}
	if len(resp.Series) != 0 {
		t.Fatalf("expected empty series, got %+v", resp.Series)
	}
}

func TestServiceLoadLinkAnalyticsUsesRequestedPeriod(t *testing.T) {
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
	service := NewService(repo, "https://tracklink.example.com")

	_, _, err := service.LoadLinkAnalytics(context.Background(), "owner-1", "link-1", LinkAnalyticsQuery{
		From:    from,
		To:      to,
		GroupBy: GroupByDay,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !capturedFrom.Equal(from) {
		t.Fatalf("expected from %s, got %s", from, capturedFrom)
	}
	if !capturedTo.Equal(to) {
		t.Fatalf("expected to %s, got %s", to, capturedTo)
	}
}

func TestServiceLoadLinkAnalyticsRejectsInvalidPeriod(t *testing.T) {
	from := time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	service := NewService(fakeDashboardRepository{}, "https://tracklink.example.com")

	_, fields, err := service.LoadLinkAnalytics(context.Background(), "owner-1", "link-1", LinkAnalyticsQuery{
		From: from,
		To:   to,
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if fields["from"] == "" {
		t.Fatalf("expected from validation field, got %v", fields)
	}
}

func TestServiceLoadRecentClicksReturnsItemsAndDefaultLimit(t *testing.T) {
	clickedAt := time.Date(2026, 5, 2, 12, 30, 0, 0, time.UTC)
	referrer := "https://t.me/example"
	userAgent := "Mozilla/5.0"
	var capturedLimit int
	repo := fakeDashboardRepository{
		getLinkByIDAndOwnerFn: func(_ context.Context, linkID, ownerID string) (links.Link, error) {
			if linkID != "link-1" {
				t.Fatalf("expected link-1, got %s", linkID)
			}
			if ownerID != "owner-1" {
				t.Fatalf("expected owner-1, got %s", ownerID)
			}
			return links.Link{ID: "link-1", OwnerID: "owner-1"}, nil
		},
		listRecentClicksFn: func(_ context.Context, linkID string, limit int) ([]ClickEvent, error) {
			if linkID != "link-1" {
				t.Fatalf("expected link-1, got %s", linkID)
			}
			capturedLimit = limit
			return []ClickEvent{
				{
					ID:        "click-1",
					LinkID:    "link-1",
					ClickedAt: clickedAt,
					Referrer:  &referrer,
					UserAgent: &userAgent,
				},
			}, nil
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	resp, fields, err := service.LoadRecentClicks(context.Background(), "owner-1", " link-1 ", RecentClicksQuery{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(fields) != 0 {
		t.Fatalf("expected empty fields, got %v", fields)
	}
	if capturedLimit != defaultRecentClicksLimit {
		t.Fatalf("expected default limit %d, got %d", defaultRecentClicksLimit, capturedLimit)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected one click, got %d", len(resp.Items))
	}
	item := resp.Items[0]
	if item.ID != "click-1" || item.LinkID != "link-1" {
		t.Fatalf("unexpected item identity: %+v", item)
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

func TestServiceLoadRecentClicksCapsLimit(t *testing.T) {
	var capturedLimit int
	repo := fakeDashboardRepository{
		listRecentClicksFn: func(_ context.Context, _ string, limit int) ([]ClickEvent, error) {
			capturedLimit = limit
			return nil, nil
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	_, _, err := service.LoadRecentClicks(context.Background(), "owner-1", "link-1", RecentClicksQuery{Limit: 150, limitProvided: true})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if capturedLimit != maxRecentClicksLimit {
		t.Fatalf("expected capped limit %d, got %d", maxRecentClicksLimit, capturedLimit)
	}
}

func TestServiceLoadRecentClicksRejectsInvalidLimit(t *testing.T) {
	tests := []struct {
		name  string
		query RecentClicksQuery
	}{
		{
			name:  "zero",
			query: RecentClicksQuery{Limit: 0, limitProvided: true},
		},
		{
			name:  "negative",
			query: RecentClicksQuery{Limit: -1, limitProvided: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(fakeDashboardRepository{}, "https://tracklink.example.com")

			_, fields, err := service.LoadRecentClicks(context.Background(), "owner-1", "link-1", tt.query)
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("expected validation error, got %v", err)
			}
			if fields["limit"] == "" {
				t.Fatalf("expected limit validation field, got %v", fields)
			}
		})
	}
}

func TestServiceLoadRecentClicksReturnsNotFoundForForeignLink(t *testing.T) {
	repo := fakeDashboardRepository{
		getLinkByIDAndOwnerFn: func(_ context.Context, _, _ string) (links.Link, error) {
			return links.Link{}, ErrLinkNotFound
		},
	}
	service := NewService(repo, "https://tracklink.example.com")

	_, _, err := service.LoadRecentClicks(context.Background(), "owner-1", "foreign-link", RecentClicksQuery{})
	if !errors.Is(err, ErrLinkNotFound) {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}
}
