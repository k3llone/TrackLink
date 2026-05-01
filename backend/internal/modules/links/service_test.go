package links

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	createFn              func(ctx context.Context, link *Link) error
	existsByCodeFn        func(ctx context.Context, code string) (bool, error)
	existsByCustomAliasFn func(ctx context.Context, customAlias string) (bool, error)
	listByOwnerFn         func(ctx context.Context, filter ListLinksFilter) ([]Link, int64, error)
	getByIDAndOwnerFn     func(ctx context.Context, linkID, ownerID string) (Link, error)
	updateStatusFn        func(ctx context.Context, linkID, ownerID, status string) (Link, error)
	softDeleteFn          func(ctx context.Context, linkID, ownerID string) error
}

func (f fakeRepository) Create(ctx context.Context, link *Link) error {
	if f.createFn == nil {
		return nil
	}

	return f.createFn(ctx, link)
}

func (f fakeRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	if f.existsByCodeFn == nil {
		return false, nil
	}

	return f.existsByCodeFn(ctx, code)
}

func (f fakeRepository) ExistsByCustomAlias(ctx context.Context, customAlias string) (bool, error) {
	if f.existsByCustomAliasFn == nil {
		return false, nil
	}

	return f.existsByCustomAliasFn(ctx, customAlias)
}

func (f fakeRepository) ListByOwner(ctx context.Context, filter ListLinksFilter) ([]Link, int64, error) {
	if f.listByOwnerFn == nil {
		return nil, 0, nil
	}

	return f.listByOwnerFn(ctx, filter)
}

func (f fakeRepository) GetByIDAndOwner(ctx context.Context, linkID, ownerID string) (Link, error) {
	if f.getByIDAndOwnerFn == nil {
		return Link{}, ErrLinkNotFound
	}

	return f.getByIDAndOwnerFn(ctx, linkID, ownerID)
}

func (f fakeRepository) UpdateStatus(ctx context.Context, linkID, ownerID, status string) (Link, error) {
	if f.updateStatusFn == nil {
		return Link{}, ErrLinkNotFound
	}

	return f.updateStatusFn(ctx, linkID, ownerID, status)
}

func (f fakeRepository) SoftDelete(ctx context.Context, linkID, ownerID string) error {
	if f.softDeleteFn == nil {
		return nil
	}

	return f.softDeleteFn(ctx, linkID, ownerID)
}

func TestServiceCreateSuccess(t *testing.T) {
	repo := fakeRepository{
		createFn: func(_ context.Context, link *Link) error {
			if link.OwnerID != "c7364dce-f6fd-4ec4-b7bc-f2a95f2f9de8" {
				t.Fatalf("unexpected owner id in create: %s", link.OwnerID)
			}
			if link.TargetURL != "https://example.com/landing" {
				t.Fatalf("unexpected target url in create: %s", link.TargetURL)
			}
			if link.Status != StatusActive {
				t.Fatalf("unexpected status in create: %s", link.Status)
			}
			link.ID = "9e1f66cf-4f5c-421f-97fa-bf4a9f34ef7a"
			return nil
		},
	}

	service := NewService(repo)
	link, fields, err := service.Create(context.Background(), "c7364dce-f6fd-4ec4-b7bc-f2a95f2f9de8", CreateLinkRequest{
		TargetURL: "https://example.com/landing",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if link.ID == "" {
		t.Fatal("expected link ID to be set")
	}
	if link.Status != StatusActive {
		t.Fatalf("expected status %q, got %q", StatusActive, link.Status)
	}
	if link.Code == "" {
		t.Fatal("expected code to be generated")
	}
}

func TestServiceCreateRetriesOnCodeCollision(t *testing.T) {
	callCount := 0
	repo := fakeRepository{
		existsByCodeFn: func(_ context.Context, code string) (bool, error) {
			return code == "COLLIDE", nil
		},
	}
	service := NewService(repo)
	service.generate = func(_ int) (string, error) {
		callCount++
		if callCount == 1 {
			return "COLLIDE", nil
		}
		return "UNIQUE1", nil
	}

	link, fields, err := service.Create(context.Background(), "owner-1", CreateLinkRequest{
		TargetURL: "https://example.com",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if link.Code != "UNIQUE1" {
		t.Fatalf("expected UNIQUE1, got %s", link.Code)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 generator calls, got %d", callCount)
	}
}

func TestServiceCreateValidation(t *testing.T) {
	repo := fakeRepository{
		createFn: func(_ context.Context, _ *Link) error {
			t.Fatal("repo must not be called on validation errors")
			return nil
		},
	}

	service := NewService(repo)
	_, fields, err := service.Create(context.Background(), "", CreateLinkRequest{
		TargetURL: " ",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if fields == nil {
		t.Fatal("expected validation fields")
	}
	if _, ok := fields["targetUrl"]; !ok {
		t.Fatal("expected targetUrl field validation error")
	}
	if _, ok := fields["ownerId"]; !ok {
		t.Fatal("expected ownerId field validation error")
	}
}

func TestServiceCreateTargetURLValidation(t *testing.T) {
	testCases := []struct {
		name      string
		targetURL string
		wantError bool
	}{
		{name: "relative URL", targetURL: "/landing", wantError: true},
		{name: "missing host", targetURL: "https://", wantError: true},
		{name: "unsupported scheme", targetURL: "ftp://example.com", wantError: true},
		{name: "valid https URL", targetURL: "https://example.com/path", wantError: false},
	}

	repo := fakeRepository{
		createFn: func(_ context.Context, _ *Link) error {
			return nil
		},
	}
	service := NewService(repo)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, fields, err := service.Create(context.Background(), "c7364dce-f6fd-4ec4-b7bc-f2a95f2f9de8", CreateLinkRequest{
				TargetURL: tc.targetURL,
			})

			if tc.wantError {
				if !errors.Is(err, ErrValidation) {
					t.Fatalf("expected ErrValidation, got %v", err)
				}
				if _, ok := fields["targetUrl"]; !ok {
					t.Fatalf("expected targetUrl field error, got %v", fields)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if fields != nil {
				t.Fatalf("expected nil fields, got %v", fields)
			}
		})
	}
}

func TestServiceCreateReturnsErrorWhenCodeGenerationAttemptsExhausted(t *testing.T) {
	repo := fakeRepository{
		existsByCodeFn: func(_ context.Context, _ string) (bool, error) {
			return true, nil
		},
	}
	service := NewService(repo)
	service.generate = func(_ int) (string, error) {
		return "always1", nil
	}

	_, _, err := service.Create(context.Background(), "owner-1", CreateLinkRequest{
		TargetURL: "https://example.com",
	})
	if !errors.Is(err, ErrCodeGenerationExhausted) {
		t.Fatalf("expected ErrCodeGenerationExhausted, got %v", err)
	}
}

func TestServiceCreateReturnsConflictWhenCustomAliasExists(t *testing.T) {
	createCalled := false
	repo := fakeRepository{
		createFn: func(_ context.Context, _ *Link) error {
			createCalled = true
			return nil
		},
		existsByCustomAliasFn: func(_ context.Context, customAlias string) (bool, error) {
			if customAlias == "taken-alias" {
				return true, nil
			}
			return false, nil
		},
	}
	service := NewService(repo)

	_, _, err := service.Create(context.Background(), "owner-1", CreateLinkRequest{
		TargetURL:   "https://example.com",
		CustomAlias: strPtr("taken-alias"),
	})
	if !errors.Is(err, ErrAliasAlreadyExists) {
		t.Fatalf("expected ErrAliasAlreadyExists, got %v", err)
	}
	if createCalled {
		t.Fatal("create should not be called when alias is already taken")
	}
}

func TestServiceCreateCustomAliasValidation(t *testing.T) {
	testCases := []struct {
		name       string
		customAlias *string
		wantError  bool
	}{
		{name: "too short", customAlias: strPtr("ab"), wantError: true},
		{name: "too long", customAlias: strPtr("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), wantError: true},
		{name: "invalid chars", customAlias: strPtr("bad alias!"), wantError: true},
		{name: "valid alias", customAlias: strPtr("spring-campaign_2026"), wantError: false},
		{name: "trimmed empty alias", customAlias: strPtr("   "), wantError: false},
	}

	repo := fakeRepository{}
	service := NewService(repo)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			link, fields, err := service.Create(context.Background(), "owner-1", CreateLinkRequest{
				TargetURL:   "https://example.com",
				CustomAlias: tc.customAlias,
			})

			if tc.wantError {
				if !errors.Is(err, ErrValidation) {
					t.Fatalf("expected ErrValidation, got %v", err)
				}
				if _, ok := fields["customAlias"]; !ok {
					t.Fatalf("expected customAlias validation error, got %v", fields)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if fields != nil {
				t.Fatalf("expected nil fields, got %v", fields)
			}
			if tc.customAlias != nil && *tc.customAlias == "   " && link.CustomAlias != nil {
				t.Fatalf("expected normalized empty alias to become nil, got %v", *link.CustomAlias)
			}
		})
	}
}

func strPtr(v string) *string {
	return &v
}

func TestServiceListSuccess(t *testing.T) {
	repo := fakeRepository{
		listByOwnerFn: func(_ context.Context, filter ListLinksFilter) ([]Link, int64, error) {
			if filter.OwnerID != "owner-1" {
				t.Fatalf("unexpected owner id: %s", filter.OwnerID)
			}
			if filter.Page != 2 {
				t.Fatalf("expected page 2, got %d", filter.Page)
			}
			if filter.PageSize != 50 {
				t.Fatalf("expected pageSize 50, got %d", filter.PageSize)
			}
			if filter.Q != "promo" {
				t.Fatalf("expected q=promo, got %s", filter.Q)
			}
			if filter.Status != StatusActive {
				t.Fatalf("expected status active, got %s", filter.Status)
			}
			return []Link{
				{
					ID:         "link-1",
					OwnerID:    "owner-1",
					Code:       "abc123",
					TargetURL:  "https://example.com",
					Status:     StatusActive,
					TotalClicks: 10,
				},
			}, 1, nil
		},
	}

	service := NewService(repo)
	items, pagination, fields, err := service.List(context.Background(), "owner-1", ListLinksQuery{
		Page:     2,
		PageSize: 50,
		Q:        "promo",
		Status:   StatusActive,
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if pagination.Page != 2 || pagination.PageSize != 50 || pagination.TotalItems != 1 || pagination.TotalPages != 1 {
		t.Fatalf("unexpected pagination: %+v", pagination)
	}
}

func TestServiceListValidation(t *testing.T) {
	service := NewService(fakeRepository{})

	_, _, fields, err := service.List(context.Background(), "", ListLinksQuery{
		Page:     0,
		PageSize: 0,
		Status:   "wrong-status",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if _, ok := fields["ownerId"]; !ok {
		t.Fatalf("expected ownerId validation error, got %v", fields)
	}
	if _, ok := fields["status"]; !ok {
		t.Fatalf("expected status validation error, got %v", fields)
	}
}

func TestServiceUpdateStatusSuccess(t *testing.T) {
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, linkID, ownerID string) (Link, error) {
			if linkID != "link-1" || ownerID != "owner-1" {
				t.Fatalf("unexpected lookup args: linkID=%s ownerID=%s", linkID, ownerID)
			}
			return Link{ID: linkID, OwnerID: ownerID, Status: StatusActive}, nil
		},
		updateStatusFn: func(_ context.Context, linkID, ownerID, status string) (Link, error) {
			if status != StatusInactive {
				t.Fatalf("expected status inactive, got %s", status)
			}
			return Link{
				ID:         linkID,
				OwnerID:    ownerID,
				Code:       "abc123",
				TargetURL:  "https://example.com",
				Status:     status,
				TotalClicks: 10,
			}, nil
		},
	}

	service := NewService(repo)
	link, fields, err := service.UpdateStatus(context.Background(), "owner-1", "link-1", UpdateLinkStatusRequest{
		Status: StatusInactive,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if link.Status != StatusInactive {
		t.Fatalf("expected inactive status, got %s", link.Status)
	}
}

func TestServiceUpdateStatusValidation(t *testing.T) {
	service := NewService(fakeRepository{})

	_, fields, err := service.UpdateStatus(context.Background(), "", "", UpdateLinkStatusRequest{
		Status: "blocked",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if _, ok := fields["ownerId"]; !ok {
		t.Fatalf("expected ownerId validation error, got %v", fields)
	}
	if _, ok := fields["linkId"]; !ok {
		t.Fatalf("expected linkId validation error, got %v", fields)
	}
	if _, ok := fields["status"]; !ok {
		t.Fatalf("expected status validation error, got %v", fields)
	}
}

func TestServiceUpdateStatusLinkNotFound(t *testing.T) {
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
			return Link{}, ErrLinkNotFound
		},
	}
	service := NewService(repo)

	_, _, err := service.UpdateStatus(context.Background(), "owner-1", "link-1", UpdateLinkStatusRequest{
		Status: StatusInactive,
	})
	if !errors.Is(err, ErrLinkNotFound) {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}
}

func TestServiceUpdateStatusForbiddenForDeletedOrBlocked(t *testing.T) {
	statuses := []string{StatusDeleted, StatusBlocked}
	for _, st := range statuses {
		st := st
		t.Run(st, func(t *testing.T) {
			repo := fakeRepository{
				getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
					return Link{ID: "link-1", OwnerID: "owner-1", Status: st}, nil
				},
				updateStatusFn: func(_ context.Context, _, _, _ string) (Link, error) {
					t.Fatal("updateStatus should not be called for blocked/deleted links")
					return Link{}, nil
				},
			}
			service := NewService(repo)

			_, _, err := service.UpdateStatus(context.Background(), "owner-1", "link-1", UpdateLinkStatusRequest{
				Status: StatusActive,
			})
			if !errors.Is(err, ErrStatusChangeNotAllowed) {
				t.Fatalf("expected ErrStatusChangeNotAllowed, got %v", err)
			}
		})
	}
}

func TestServiceDeleteSuccess(t *testing.T) {
	softDeleted := false
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, linkID, ownerID string) (Link, error) {
			return Link{ID: linkID, OwnerID: ownerID, Status: StatusActive}, nil
		},
		softDeleteFn: func(_ context.Context, linkID, ownerID string) error {
			softDeleted = true
			if linkID != "link-1" || ownerID != "owner-1" {
				t.Fatalf("unexpected delete args: %s %s", linkID, ownerID)
			}
			return nil
		},
	}
	service := NewService(repo)

	fields, err := service.Delete(context.Background(), "owner-1", "link-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if !softDeleted {
		t.Fatal("expected soft delete to be called")
	}
}

func TestServiceDeleteValidation(t *testing.T) {
	service := NewService(fakeRepository{})

	fields, err := service.Delete(context.Background(), "", "")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if _, ok := fields["ownerId"]; !ok {
		t.Fatalf("expected ownerId validation error, got %v", fields)
	}
	if _, ok := fields["linkId"]; !ok {
		t.Fatalf("expected linkId validation error, got %v", fields)
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
			return Link{}, ErrLinkNotFound
		},
	}
	service := NewService(repo)

	_, err := service.Delete(context.Background(), "owner-1", "missing")
	if !errors.Is(err, ErrLinkNotFound) {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}
}

func TestServiceDeleteAlreadyDeletedIsIdempotent(t *testing.T) {
	softDeleteCalled := false
	now := timePtr(time.Now().UTC())
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
			return Link{
				ID:        "link-1",
				OwnerID:   "owner-1",
				Status:    StatusDeleted,
				DeletedAt: now,
			}, nil
		},
		softDeleteFn: func(_ context.Context, _, _ string) error {
			softDeleteCalled = true
			return nil
		},
	}
	service := NewService(repo)

	fields, err := service.Delete(context.Background(), "owner-1", "link-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if softDeleteCalled {
		t.Fatal("soft delete should not be called for already deleted link")
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
