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
	if len(steps) != 2 || steps[0] != "touch" || steps[1] != "event" {
		t.Fatalf("expected touch before click event, got %v", steps)
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

func TestServiceResolveUnavailableDoesNotTrackOrTouch(t *testing.T) {
	deletedAt := time.Date(2026, 5, 2, 11, 0, 0, 0, time.UTC)
	tests := []struct {
		name        string
		link        Link
		wantStatus  string
		wantDeleted bool
	}{
		{
			name: "blocked",
			link: Link{
				ID:        "link-blocked",
				Code:      "blocked-link",
				TargetURL: "https://example.com/blocked",
				Status:    StatusBlocked,
			},
			wantStatus: StatusBlocked,
		},
		{
			name: "inactive",
			link: Link{
				ID:        "link-inactive",
				Code:      "inactive-link",
				TargetURL: "https://example.com/inactive",
				Status:    StatusInactive,
			},
			wantStatus: StatusInactive,
		},
		{
			name: "deleted status",
			link: Link{
				ID:        "link-deleted",
				Code:      "deleted-link",
				TargetURL: "https://example.com/deleted",
				Status:    StatusDeleted,
			},
			wantStatus:  StatusDeleted,
			wantDeleted: true,
		},
		{
			name: "soft deleted",
			link: Link{
				ID:        "link-soft-deleted",
				Code:      "soft-deleted-link",
				TargetURL: "https://example.com/soft-deleted",
				Status:    StatusActive,
				DeletedAt: &deletedAt,
			},
			wantStatus:  StatusActive,
			wantDeleted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyticsRepo := fakeAnalyticsRepository{
				createClickEventFn: func(_ context.Context, _ analytics.CreateClickEventParams) error {
					t.Fatal("create click event must not be called for unavailable link")
					return nil
				},
			}
			repo := fakeRepository{
				findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
					if code != tt.link.Code {
						t.Fatalf("expected code %s, got %s", tt.link.Code, code)
					}
					return tt.link, nil
				},
				touchActiveLinkFn: func(_ context.Context, _ string, _ time.Time) error {
					t.Fatal("touch active link must not be called for unavailable link")
					return nil
				},
			}

			service := NewService(repo, analyticsRepo)
			result, err := service.ResolveAndTrack(context.Background(), tt.link.Code, RequestMeta{})
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if result.Kind != ResultKindUnavailable {
				t.Fatalf("expected unavailable result, got %s", result.Kind)
			}
			if result.Status != tt.wantStatus {
				t.Fatalf("expected status %s, got %s", tt.wantStatus, result.Status)
			}
			if result.Deleted != tt.wantDeleted {
				t.Fatalf("expected deleted flag %v, got %v", tt.wantDeleted, result.Deleted)
			}
		})
	}
}

func TestServiceResolveDoesNotTrackWhenActiveTouchFindsNoActiveRow(t *testing.T) {
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, _ string) (Link, error) {
			return Link{
				ID:        "link-race",
				Code:      "race",
				TargetURL: "https://example.com/race",
				Status:    StatusActive,
			}, nil
		},
		touchActiveLinkFn: func(_ context.Context, _ string, _ time.Time) error {
			return ErrLinkNotFound
		},
	}
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, _ analytics.CreateClickEventParams) error {
			t.Fatal("create click event must not be called when active touch fails")
			return nil
		},
	}

	service := NewService(repo, analyticsRepo)
	result, err := service.ResolveAndTrack(context.Background(), "race", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Kind != ResultKindUnavailable {
		t.Fatalf("expected unavailable result, got %s", result.Kind)
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

func TestServiceResolveNormalizesClickedAtToUTC(t *testing.T) {
	var capturedClickedAt time.Time
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, _ string) (Link, error) {
			return Link{
				ID:        "link-utc",
				Code:      "tz",
				TargetURL: "https://example.com/tz",
				Status:    StatusActive,
			}, nil
		},
	}
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, event analytics.CreateClickEventParams) error {
			capturedClickedAt = event.ClickedAt
			return nil
		},
	}

	localClickedAt := time.Date(2026, 5, 2, 16, 45, 0, 0, time.FixedZone("UTC+3", 3*60*60))
	service := NewService(repo, analyticsRepo)
	_, err := service.ResolveAndTrack(context.Background(), "tz", RequestMeta{
		ClickedAt: localClickedAt,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if capturedClickedAt.IsZero() {
		t.Fatal("expected analytics repository to receive clicked at value")
	}
	if capturedClickedAt.Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", capturedClickedAt.Location())
	}
	if !capturedClickedAt.Equal(localClickedAt.UTC()) {
		t.Fatalf("expected clickedAt %v, got %v", localClickedAt.UTC(), capturedClickedAt)
	}
}

func TestServiceResolveSetsClickedAtWhenMetaIsZero(t *testing.T) {
	var capturedClickedAt time.Time
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, _ string) (Link, error) {
			return Link{
				ID:        "link-now",
				Code:      "now",
				TargetURL: "https://example.com/now",
				Status:    StatusActive,
			}, nil
		},
	}
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, event analytics.CreateClickEventParams) error {
			capturedClickedAt = event.ClickedAt
			return nil
		},
	}

	service := NewService(repo, analyticsRepo)
	_, err := service.ResolveAndTrack(context.Background(), "now", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if capturedClickedAt.IsZero() {
		t.Fatal("expected non-zero clickedAt when request meta is zero")
	}
	if capturedClickedAt.Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", capturedClickedAt.Location())
	}
}

func TestServiceResolveTracksClickWithResolvedLinkIDForAlias(t *testing.T) {
	var capturedLinkID string
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != "campaign" {
				t.Fatalf("expected alias campaign, got %s", code)
			}
			return Link{
				ID:          "link-alias-42",
				Code:        "abc123",
				CustomAlias: ptrString("campaign"),
				TargetURL:   "https://example.com/campaign",
				Status:      StatusActive,
			}, nil
		},
	}
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, event analytics.CreateClickEventParams) error {
			capturedLinkID = event.LinkID
			return nil
		},
	}

	service := NewService(repo, analyticsRepo)
	_, err := service.ResolveAndTrack(context.Background(), "campaign", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if capturedLinkID != "link-alias-42" {
		t.Fatalf("expected click event linked to link-alias-42, got %s", capturedLinkID)
	}
}

func TestServiceResolveNotFoundDoesNotCreateClickEvent(t *testing.T) {
	wasCalled := false
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, _ string) (Link, error) {
			return Link{}, ErrLinkNotFound
		},
	}
	analyticsRepo := fakeAnalyticsRepository{
		createClickEventFn: func(_ context.Context, _ analytics.CreateClickEventParams) error {
			wasCalled = true
			return nil
		},
	}

	service := NewService(repo, analyticsRepo)
	result, err := service.ResolveAndTrack(context.Background(), "missing", RequestMeta{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.Kind != ResultKindNotFound {
		t.Fatalf("expected not found result kind, got %s", result.Kind)
	}
	if wasCalled {
		t.Fatal("expected click event not to be created for missing link")
	}
}

func ptrString(value string) *string {
	return &value
}
