package links

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

func TestHandlerCreateSuccess(t *testing.T) {
	repo := fakeRepository{
		createFn: func(_ context.Context, link *Link) error {
			link.ID = "7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a"
			now := time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC)
			link.CreatedAt = now
			link.UpdatedAt = now
			return nil
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, sessionID string) (session.SessionData, error) {
			if sessionID != "session-1" {
				return session.SessionData{}, errors.New("unknown session")
			}
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", strings.NewReader(`{"targetUrl":"https://example.com/page"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()

	auth.RequireAuth(http.HandlerFunc(handler.Create)).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp LinkResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.OwnerID != "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a" {
		t.Fatalf("unexpected owner id: %s", resp.OwnerID)
	}
	if !strings.HasPrefix(resp.ShortURL, "https://tracklink.example.com/s/") {
		t.Fatalf("unexpected shortUrl: %s", resp.ShortURL)
	}
}

func TestHandlerCreateUnauthorizedWithoutSessionData(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", strings.NewReader(`{"targetUrl":"https://example.com"}`))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestHandlerCreateInvalidTargetURL(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, sessionID string) (session.SessionData, error) {
			if sessionID != "session-1" {
				return session.SessionData{}, errors.New("unknown session")
			}
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", strings.NewReader(`{"targetUrl":"not-an-url"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	auth.RequireAuth(http.HandlerFunc(handler.Create)).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "validation_error" {
		t.Fatalf("unexpected error code: %s", resp.Error.Code)
	}
	if _, ok := resp.Error.Fields["targetUrl"]; !ok {
		t.Fatalf("expected targetUrl field error, got %v", resp.Error.Fields)
	}
}

func TestHandlerCreateWithCustomAliasUsesAliasInShortURL(t *testing.T) {
	repo := fakeRepository{
		createFn: func(_ context.Context, link *Link) error {
			link.ID = "7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a"
			now := time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC)
			link.CreatedAt = now
			link.UpdatedAt = now
			return nil
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, sessionID string) (session.SessionData, error) {
			if sessionID != "session-1" {
				return session.SessionData{}, errors.New("unknown session")
			}
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", strings.NewReader(`{"targetUrl":"https://example.com/page","customAlias":"spring-campaign"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	auth.RequireAuth(http.HandlerFunc(handler.Create)).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp LinkResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.CustomAlias == nil || *resp.CustomAlias != "spring-campaign" {
		t.Fatalf("expected customAlias spring-campaign, got %v", resp.CustomAlias)
	}
	if resp.ShortURL != "https://tracklink.example.com/s/spring-campaign" {
		t.Fatalf("unexpected shortUrl: %s", resp.ShortURL)
	}
}

func TestHandlerCreateReturnsConflictForTakenAlias(t *testing.T) {
	repo := fakeRepository{
		existsByCustomAliasFn: func(_ context.Context, customAlias string) (bool, error) {
			return customAlias == "taken-alias", nil
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, sessionID string) (session.SessionData, error) {
			if sessionID != "session-1" {
				return session.SessionData{}, errors.New("unknown session")
			}
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", strings.NewReader(`{"targetUrl":"https://example.com/page","customAlias":"taken-alias"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	auth.RequireAuth(http.HandlerFunc(handler.Create)).ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "custom_alias_already_exists" {
		t.Fatalf("unexpected error code: %s", resp.Error.Code)
	}
}

func TestHandlerListSuccessWithQueryParams(t *testing.T) {
	alias := "spring"
	repo := fakeRepository{
		listByOwnerFn: func(_ context.Context, filter ListLinksFilter) ([]Link, int64, error) {
			if filter.OwnerID != "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a" {
				t.Fatalf("unexpected ownerID: %s", filter.OwnerID)
			}
			if filter.Page != 2 {
				t.Fatalf("expected page=2, got %d", filter.Page)
			}
			if filter.PageSize != 30 {
				t.Fatalf("expected pageSize=30, got %d", filter.PageSize)
			}
			if filter.Q != "promo" {
				t.Fatalf("expected q=promo, got %s", filter.Q)
			}
			if filter.Status != StatusActive {
				t.Fatalf("expected status=active, got %s", filter.Status)
			}

			return []Link{
				{
					ID:          "7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a",
					OwnerID:     "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
					Code:        "abc123",
					CustomAlias: &alias,
					TargetURL:   "https://example.com/page",
					Status:      StatusActive,
					TotalClicks: 42,
					CreatedAt:   time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
				},
			}, 1, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, sessionID string) (session.SessionData, error) {
			if sessionID != "session-1" {
				return session.SessionData{}, errors.New("unknown session")
			}
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links?page=2&pageSize=30&q=promo&status=active", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()

	auth.RequireAuth(http.HandlerFunc(handler.List)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp LinkListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Pagination.Page != 2 || resp.Pagination.PageSize != 30 || resp.Pagination.TotalItems != 1 || resp.Pagination.TotalPages != 1 {
		t.Fatalf("unexpected pagination: %+v", resp.Pagination)
	}
}

func TestHandlerListSearchQueryByCodeAliasAndTargetURL(t *testing.T) {
	testCases := []struct {
		name string
		q    string
	}{
		{name: "code or alias search", q: "spring"},
		{name: "target url search", q: "example.com/landing"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := fakeRepository{
				listByOwnerFn: func(_ context.Context, filter ListLinksFilter) ([]Link, int64, error) {
					if filter.OwnerID != "owner-1" {
						t.Fatalf("unexpected ownerID: %s", filter.OwnerID)
					}
					if filter.Q != tc.q {
						t.Fatalf("expected q=%s, got %s", tc.q, filter.Q)
					}
					return []Link{}, 0, nil
				},
			}

			handler := NewHandler(NewService(repo), "https://tracklink.example.com")
			auth := httpmiddleware.NewAuth(fakeSessionStore{
				getFn: func(_ context.Context, _ string) (session.SessionData, error) {
					return session.SessionData{
						UserID: "owner-1",
						Role:   "customer",
					}, nil
				},
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/links?q="+tc.q, nil)
			req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
			rr := httptest.NewRecorder()

			auth.RequireAuth(http.HandlerFunc(handler.List)).ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHandlerListInvalidQueryReturnsValidationError(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links?page=abc&pageSize=-1", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()

	auth.RequireAuth(http.HandlerFunc(handler.List)).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "validation_error" {
		t.Fatalf("unexpected error code: %s", resp.Error.Code)
	}
}

func TestHandlerListRejectsPageSizeAboveMax(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links?pageSize=101", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()

	auth.RequireAuth(http.HandlerFunc(handler.List)).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := resp.Error.Fields["pageSize"]; !ok {
		t.Fatalf("expected pageSize field error, got %v", resp.Error.Fields)
	}
}

func TestHandlerListEmptyResultReturnsEmptyItemsArray(t *testing.T) {
	repo := fakeRepository{
		listByOwnerFn: func(_ context.Context, filter ListLinksFilter) ([]Link, int64, error) {
			if filter.OwnerID != "owner-1" {
				t.Fatalf("unexpected ownerID: %s", filter.OwnerID)
			}
			return nil, 0, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{
				UserID: "owner-1",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links?q=missing", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()

	auth.RequireAuth(http.HandlerFunc(handler.List)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	itemsRaw, ok := payload["items"]
	if !ok {
		t.Fatalf("items key is missing in response: %v", payload)
	}
	items, ok := itemsRaw.([]any)
	if !ok {
		t.Fatalf("items must be an array, got %T", itemsRaw)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty items array, got len=%d", len(items))
	}
}

func TestHandlerListResponseContainsFullLinkFields(t *testing.T) {
	alias := "campaign"
	lastClicked := time.Date(2026, 5, 1, 9, 30, 0, 0, time.UTC)
	repo := fakeRepository{
		listByOwnerFn: func(_ context.Context, _ ListLinksFilter) ([]Link, int64, error) {
			return []Link{
				{
					ID:            "c37ca8d2-39aa-44a0-a432-a06bf92f19a4",
					OwnerID:       "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
					Code:          "abc123",
					CustomAlias:   &alias,
					TargetURL:     "https://example.com/landing",
					Status:        StatusActive,
					TotalClicks:   42,
					LastClickedAt: &lastClicked,
					CreatedAt:     time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC),
					UpdatedAt:     time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
				},
				{
					ID:          "1a15af7a-f9be-44e3-997e-c5d15f4ac32a",
					OwnerID:     "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
					Code:        "noalias",
					TargetURL:   "https://example.com/noalias",
					Status:      StatusInactive,
					TotalClicks: 0,
					CreatedAt:   time.Date(2026, 4, 30, 7, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2026, 4, 30, 8, 0, 0, 0, time.UTC),
				},
			}, 2, nil
		},
	}

	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{
				UserID: "b3f4c113-6f22-42f2-8b45-6e88b2f9b71a",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	auth.RequireAuth(http.HandlerFunc(handler.List)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp LinkListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	first := resp.Items[0]
	if first.ID == "" || first.OwnerID == "" || first.Code == "" || first.ShortURL == "" || first.TargetURL == "" || first.Status == "" || first.CreatedAt == "" || first.UpdatedAt == "" {
		t.Fatalf("first item has missing required fields: %+v", first)
	}
	if first.CustomAlias == nil || *first.CustomAlias != "campaign" {
		t.Fatalf("expected custom alias campaign, got %v", first.CustomAlias)
	}
	if first.ShortURL != "https://tracklink.example.com/s/campaign" {
		t.Fatalf("unexpected shortUrl for alias item: %s", first.ShortURL)
	}
	if first.LastClickedAt == nil {
		t.Fatal("expected lastClickedAt to be present")
	}

	second := resp.Items[1]
	if second.CustomAlias != nil {
		t.Fatalf("expected nil customAlias for second item, got %v", second.CustomAlias)
	}
	if second.ShortURL != "https://tracklink.example.com/s/noalias" {
		t.Fatalf("unexpected shortUrl for code item: %s", second.ShortURL)
	}
	if second.LastClickedAt != nil {
		t.Fatalf("expected nil lastClickedAt for second item, got %v", second.LastClickedAt)
	}
}

func TestHandlerUpdateStatusSuccess(t *testing.T) {
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, linkID, ownerID string) (Link, error) {
			return Link{ID: linkID, OwnerID: ownerID, Status: StatusActive}, nil
		},
		updateStatusFn: func(_ context.Context, linkID, ownerID, status string) (Link, error) {
			return Link{
				ID:        linkID,
				OwnerID:   ownerID,
				Code:      "abc123",
				TargetURL: "https://example.com",
				Status:    status,
				CreatedAt: time.Date(2026, 5, 1, 7, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
			}, nil
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{
				UserID: "owner-1",
				Role:   "customer",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/links/"+testLinkID+"/status", strings.NewReader(`{"status":"inactive"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Patch("/api/v1/links/{linkId}/status", handler.UpdateStatus)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp LinkResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != StatusInactive {
		t.Fatalf("expected inactive status, got %s", resp.Status)
	}
}

func TestHandlerUpdateStatusInvalidBody(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/links/"+testLinkID+"/status", strings.NewReader(`{"status"`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Patch("/api/v1/links/{linkId}/status", handler.UpdateStatus)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandlerUpdateStatusInvalidLinkID(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/links/not-a-uuid/status", strings.NewReader(`{"status":"active"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Patch("/api/v1/links/{linkId}/status", handler.UpdateStatus)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := resp.Error.Fields["linkId"]; !ok {
		t.Fatalf("expected linkId field error, got %v", resp.Error.Fields)
	}
}

func TestHandlerUpdateStatusNotFound(t *testing.T) {
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
			return Link{}, ErrLinkNotFound
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/links/"+testLinkID+"/status", strings.NewReader(`{"status":"active"}`))
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Patch("/api/v1/links/{linkId}/status", handler.UpdateStatus)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandlerUpdateStatusConflictForBlockedOrDeleted(t *testing.T) {
	for _, st := range []string{StatusBlocked, StatusDeleted} {
		st := st
		t.Run(st, func(t *testing.T) {
			repo := fakeRepository{
				getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
					return Link{ID: testLinkID, OwnerID: "owner-1", Status: st}, nil
				},
			}
			handler := NewHandler(NewService(repo), "https://tracklink.example.com")
			auth := httpmiddleware.NewAuth(fakeSessionStore{
				getFn: func(_ context.Context, _ string) (session.SessionData, error) {
					return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
				},
			})

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/links/"+testLinkID+"/status", strings.NewReader(`{"status":"active"}`))
			req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
			rr := httptest.NewRecorder()
			testRouter := chi.NewRouter()
			testRouter.With(auth.RequireAuth).Patch("/api/v1/links/{linkId}/status", handler.UpdateStatus)
			testRouter.ServeHTTP(rr, req)

			if rr.Code != http.StatusConflict {
				t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
			}
		})
	}
}

func TestHandlerDeleteSuccess(t *testing.T) {
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, linkID, ownerID string) (Link, error) {
			return Link{ID: linkID, OwnerID: ownerID, Status: StatusActive}, nil
		},
		softDeleteFn: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/"+testLinkID, nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Delete("/api/v1/links/{linkId}", handler.Delete)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "" {
		t.Fatalf("expected empty body for 204, got %q", rr.Body.String())
	}
}

func TestHandlerDeleteInvalidLinkID(t *testing.T) {
	handler := NewHandler(NewService(fakeRepository{}), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/not-a-uuid", nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Delete("/api/v1/links/{linkId}", handler.Delete)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := resp.Error.Fields["linkId"]; !ok {
		t.Fatalf("expected linkId field error, got %v", resp.Error.Fields)
	}
}

func TestHandlerDeleteNotFound(t *testing.T) {
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
			return Link{}, ErrLinkNotFound
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/"+testMissingLinkID, nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Delete("/api/v1/links/{linkId}", handler.Delete)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandlerDeleteIdempotentForAlreadyDeletedLink(t *testing.T) {
	now := time.Now().UTC()
	repo := fakeRepository{
		getByIDAndOwnerFn: func(_ context.Context, _, _ string) (Link, error) {
			return Link{
				ID:        testLinkID,
				OwnerID:   "owner-1",
				Status:    StatusDeleted,
				DeletedAt: &now,
			}, nil
		},
		softDeleteFn: func(_ context.Context, _, _ string) error {
			t.Fatal("softDelete should not be called for already deleted link")
			return nil
		},
	}
	handler := NewHandler(NewService(repo), "https://tracklink.example.com")
	auth := httpmiddleware.NewAuth(fakeSessionStore{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{UserID: "owner-1", Role: "customer"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/"+testLinkID, nil)
	req.AddCookie(&http.Cookie{Name: httpmiddleware.SessionCookieName, Value: "session-1"})
	rr := httptest.NewRecorder()
	testRouter := chi.NewRouter()
	testRouter.With(auth.RequireAuth).Delete("/api/v1/links/{linkId}", handler.Delete)
	testRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}
