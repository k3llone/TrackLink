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
	if !strings.HasPrefix(resp.ShortURL, "https://tracklink.example.com/") {
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
	if resp.ShortURL != "https://tracklink.example.com/spring-campaign" {
		t.Fatalf("unexpected shortUrl: %s", resp.ShortURL)
	}
}

