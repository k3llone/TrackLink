package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	httpmiddleware "tracklink/internal/http/middleware"
	"tracklink/internal/modules/accounts"
	"tracklink/internal/modules/links"
	"tracklink/internal/platform/session"
)

type fakeSessionStore struct {
	getFn func(ctx context.Context, sessionID string) (session.SessionData, error)
}

func (f fakeSessionStore) Get(ctx context.Context, sessionID string) (session.SessionData, error) {
	if f.getFn == nil {
		return session.SessionData{}, session.ErrSessionNotFound
	}

	return f.getFn(ctx, sessionID)
}

func TestHandlerListLinksSuccessForAdmin(t *testing.T) {
	repo := fakeRepository{
		listFn: func(_ context.Context, filter ListLinksFilter) ([]links.Link, int64, error) {
			if filter.Page != 2 {
				t.Fatalf("expected page 2, got %d", filter.Page)
			}
			if filter.PageSize != 30 {
				t.Fatalf("expected pageSize 30, got %d", filter.PageSize)
			}
			if filter.Q != "spring" {
				t.Fatalf("expected q=spring, got %s", filter.Q)
			}

			alias := "spring"
			return []links.Link{
				{
					ID:            "link-1",
					OwnerID:       "owner-1",
					Code:          "abc123",
					CustomAlias:   &alias,
					TargetURL:     "https://example.com/page",
					Status:        links.StatusActive,
					TotalClicks:   42,
					LastClickedAt: timePtr(time.Date(2026, 5, 2, 9, 0, 0, 0, time.UTC)),
					CreatedAt:     time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC),
					UpdatedAt:     time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
				},
			}, 1, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/links?page=2&pageSize=30&q=spring", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp AdminLinkListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Pagination.Page != 2 || resp.Pagination.PageSize != 30 || resp.Pagination.TotalItems != 1 || resp.Pagination.TotalPages != 1 {
		t.Fatalf("unexpected pagination: %+v", resp.Pagination)
	}
	if resp.Items[0].Status != links.StatusActive {
		t.Fatalf("expected active status, got %s", resp.Items[0].Status)
	}
}

func TestHandlerListLinksInvalidQueryReturnsValidationError(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/links?page=abc&pageSize=-1", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandlerBlockLinkSuccessForAdmin(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, linkID string) (links.Link, error) {
			return links.Link{ID: linkID, Status: links.StatusActive}, nil
		},
		updateStatusFn: func(_ context.Context, linkID, status string) (links.Link, error) {
			return links.Link{
				ID:        linkID,
				OwnerID:   "owner-1",
				Code:      "abc123",
				TargetURL: "https://example.com",
				Status:    status,
				CreatedAt: time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
			}, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/links/link-1/block", strings.NewReader(`{"reason":"spam"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp AdminLink
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != links.StatusBlocked {
		t.Fatalf("expected blocked status, got %s", resp.Status)
	}
}

func TestHandlerBlockLinkAcceptsEmptyBody(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, linkID string) (links.Link, error) {
			return links.Link{ID: linkID, Status: links.StatusBlocked}, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/links/link-1/block", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestHandlerBlockLinkInvalidBodyReturnsValidationError(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/links/link-1/block", strings.NewReader(`{"reason"`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandlerBlockLinkReturnsNotFound(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, _ string) (links.Link, error) {
			return links.Link{}, ErrLinkNotFound
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/links/missing/block", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandlerDeactivateLinkSuccessForAdmin(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, linkID string) (links.Link, error) {
			return links.Link{ID: linkID, Status: links.StatusActive}, nil
		},
		updateStatusFn: func(_ context.Context, linkID, status string) (links.Link, error) {
			if status != links.StatusInactive {
				t.Fatalf("expected inactive status update, got %s", status)
			}
			return links.Link{
				ID:        linkID,
				OwnerID:   "owner-1",
				Code:      "abc123",
				TargetURL: "https://example.com",
				Status:    status,
				CreatedAt: time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
			}, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/links/link-1/deactivate", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp AdminLink
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != links.StatusInactive {
		t.Fatalf("expected inactive status, got %s", resp.Status)
	}
}

func TestHandlerDeactivateBlockedLinkReturnsConflict(t *testing.T) {
	repo := fakeRepository{
		getByIDFn: func(_ context.Context, linkID string) (links.Link, error) {
			return links.Link{ID: linkID, Status: links.StatusBlocked}, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	router := newAdminTestRouter(handler, accounts.RoleAdmin)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/links/link-1/deactivate", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-admin"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
	}
}

func TestHandlerAdminEndpointsRequireAdmin(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		sessionID  string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{
			name:       "no session on list",
			method:     http.MethodGet,
			path:       "/api/v1/admin/links",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "customer on list",
			role:       accounts.RoleCustomer,
			sessionID:  "session-customer",
			method:     http.MethodGet,
			path:       "/api/v1/admin/links",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "no session on block",
			method:     http.MethodPatch,
			path:       "/api/v1/admin/links/link-1/block",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "customer on block",
			role:       accounts.RoleCustomer,
			sessionID:  "session-customer",
			method:     http.MethodPatch,
			path:       "/api/v1/admin/links/link-1/block",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "no session on deactivate",
			method:     http.MethodPatch,
			path:       "/api/v1/admin/links/link-1/deactivate",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "customer on deactivate",
			role:       accounts.RoleCustomer,
			sessionID:  "session-customer",
			method:     http.MethodPatch,
			path:       "/api/v1/admin/links/link-1/deactivate",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
			router := newAdminTestRouter(handler, tt.role)

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: tt.sessionID})
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}
		})
	}
}

func newAdminTestRouter(handler *Handler, role string) http.Handler {
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, sessionID string) (session.SessionData, error) {
			switch sessionID {
			case "session-admin":
				return session.SessionData{
					UserID: "admin-1",
					Role:   accounts.RoleAdmin,
				}, nil
			case "session-customer":
				return session.SessionData{
					UserID: "customer-1",
					Role:   role,
				}, nil
			default:
				return session.SessionData{}, errors.New("unknown session")
			}
		},
	})

	router := chi.NewRouter()
	router.With(auth.RequireAuth, auth.RequireAdmin).Get("/api/v1/admin/links", handler.ListLinks)
	router.With(auth.RequireAuth, auth.RequireAdmin).Patch("/api/v1/admin/links/{linkId}/block", handler.BlockLink)
	router.With(auth.RequireAuth, auth.RequireAdmin).Patch("/api/v1/admin/links/{linkId}/deactivate", handler.DeactivateLink)

	return router
}
