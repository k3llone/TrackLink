package admin

import (
	"context"
	"errors"
	"testing"
	"time"

	"tracklink/internal/modules/links"
)

type fakeRepository struct {
	listFn         func(ctx context.Context, filter ListLinksFilter) ([]links.Link, int64, error)
	getByIDFn      func(ctx context.Context, linkID string) (links.Link, error)
	updateStatusFn func(ctx context.Context, linkID, status string) (links.Link, error)
}

func (f fakeRepository) List(ctx context.Context, filter ListLinksFilter) ([]links.Link, int64, error) {
	if f.listFn == nil {
		return nil, 0, nil
	}

	return f.listFn(ctx, filter)
}

func (f fakeRepository) GetByID(ctx context.Context, linkID string) (links.Link, error) {
	if f.getByIDFn == nil {
		return links.Link{}, ErrLinkNotFound
	}

	return f.getByIDFn(ctx, linkID)
}

func (f fakeRepository) UpdateStatus(ctx context.Context, linkID, status string) (links.Link, error) {
	if f.updateStatusFn == nil {
		return links.Link{}, ErrLinkNotFound
	}

	return f.updateStatusFn(ctx, linkID, status)
}

func TestServiceListSuccess(t *testing.T) {
	repo := fakeRepository{
		listFn: func(_ context.Context, filter ListLinksFilter) ([]links.Link, int64, error) {
			if filter.Page != 2 {
				t.Fatalf("expected page 2, got %d", filter.Page)
			}
			if filter.PageSize != 50 {
				t.Fatalf("expected pageSize 50, got %d", filter.PageSize)
			}
			return []links.Link{
				{
					ID:        "link-1",
					OwnerID:   "owner-1",
					Code:      "abc123",
					TargetURL: "https://example.com",
					Status:    links.StatusActive,
				},
			}, 1, nil
		},
	}

	service := NewService(repo)
	items, pagination, fields, err := service.List(context.Background(), "admin-1", ListLinksQuery{
		Page:     2,
		PageSize: 50,
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

func TestServiceListUsesDefaultPagination(t *testing.T) {
	repo := fakeRepository{
		listFn: func(_ context.Context, filter ListLinksFilter) ([]links.Link, int64, error) {
			if filter.Page != defaultListPage {
				t.Fatalf("expected default page %d, got %d", defaultListPage, filter.Page)
			}
			if filter.PageSize != defaultListPageSize {
				t.Fatalf("expected default pageSize %d, got %d", defaultListPageSize, filter.PageSize)
			}
			return nil, 0, nil
		},
	}

	service := NewService(repo)
	_, pagination, fields, err := service.List(context.Background(), "admin-1", ListLinksQuery{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if pagination.Page != defaultListPage || pagination.PageSize != defaultListPageSize {
		t.Fatalf("unexpected pagination: %+v", pagination)
	}
}

func TestServiceBlockTransitionsActiveToBlocked(t *testing.T) {
	updateCalled := false
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, linkID string) (links.Link, error) {
			if linkID != "link-1" {
				t.Fatalf("unexpected link id: %s", linkID)
			}
			return links.Link{ID: linkID, Status: links.StatusActive}, nil
		},
		updateStatusFn: func(_ context.Context, linkID, status string) (links.Link, error) {
			updateCalled = true
			if status != links.StatusBlocked {
				t.Fatalf("expected blocked status, got %s", status)
			}
			return links.Link{ID: linkID, Status: status}, nil
		},
	}

	service := NewService(repo)
	link, fields, err := service.Block(context.Background(), "admin-1", "link-1", AdminBlockLinkRequest{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if !updateCalled {
		t.Fatal("expected update status to be called")
	}
	if link.Status != links.StatusBlocked {
		t.Fatalf("expected blocked status, got %s", link.Status)
	}
}

func TestServiceBlockTransitionsInactiveToBlocked(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, _ string) (links.Link, error) {
			return links.Link{ID: "link-1", Status: links.StatusInactive}, nil
		},
		updateStatusFn: func(_ context.Context, linkID, status string) (links.Link, error) {
			return links.Link{ID: linkID, Status: status}, nil
		},
	}

	service := NewService(repo)
	link, _, err := service.Block(context.Background(), "admin-1", "link-1", AdminBlockLinkRequest{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if link.Status != links.StatusBlocked {
		t.Fatalf("expected blocked status, got %s", link.Status)
	}
}

func TestServiceBlockReturnsBlockedLinkAsIsWhenAlreadyBlocked(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, _ string) (links.Link, error) {
			return links.Link{ID: "link-1", Status: links.StatusBlocked}, nil
		},
		updateStatusFn: func(_ context.Context, _, _ string) (links.Link, error) {
			t.Fatal("update status must not be called for blocked link")
			return links.Link{}, nil
		},
	}

	service := NewService(repo)
	link, fields, err := service.Block(context.Background(), "admin-1", "link-1", AdminBlockLinkRequest{})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if link.Status != links.StatusBlocked {
		t.Fatalf("expected blocked status, got %s", link.Status)
	}
}

func TestServiceBlockReturnsNotFoundWhenLinkMissingOrDeleted(t *testing.T) {
	tests := []struct {
		name string
		repo fakeRepository
	}{
		{
			name: "missing",
			repo: fakeRepository{
				getByIDFn: func(_ context.Context, _ string) (links.Link, error) {
					return links.Link{}, ErrLinkNotFound
				},
			},
		},
		{
			name: "deleted",
			repo: fakeRepository{
				getByIDFn: func(_ context.Context, _ string) (links.Link, error) {
					now := time.Now().UTC()
					return links.Link{
						ID:        "link-1",
						Status:    links.StatusDeleted,
						DeletedAt: &now,
					}, nil
				},
				updateStatusFn: func(_ context.Context, _, _ string) (links.Link, error) {
					t.Fatal("update status must not be called for deleted link")
					return links.Link{}, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.repo)
			_, _, err := service.Block(context.Background(), "admin-1", "link-1", AdminBlockLinkRequest{})
			if !errors.Is(err, ErrLinkNotFound) {
				t.Fatalf("expected ErrLinkNotFound, got %v", err)
			}
		})
	}
}

func TestServiceBlockTrimsReasonWithoutAffectingStorage(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, _ string) (links.Link, error) {
			return links.Link{ID: "link-1", Status: links.StatusActive}, nil
		},
		updateStatusFn: func(_ context.Context, linkID, status string) (links.Link, error) {
			return links.Link{ID: linkID, Status: status}, nil
		},
	}

	service := NewService(repo)
	link, fields, err := service.Block(context.Background(), "admin-1", "link-1", AdminBlockLinkRequest{
		Reason: strPtr("  spam  "),
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil fields, got %v", fields)
	}
	if link.Status != links.StatusBlocked {
		t.Fatalf("expected blocked status, got %s", link.Status)
	}
}

func strPtr(value string) *string {
	return &value
}
