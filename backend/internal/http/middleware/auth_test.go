package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tracklink/internal/platform/session"
)

type fakeSessionReader struct {
	getFn func(ctx context.Context, sessionID string) (session.SessionData, error)
}

func (f fakeSessionReader) Get(ctx context.Context, sessionID string) (session.SessionData, error) {
	if f.getFn == nil {
		return session.SessionData{}, session.ErrSessionNotFound
	}
	return f.getFn(ctx, sessionID)
}

func TestRequireAuthWithoutCookieReturnsUnauthorized(t *testing.T) {
	auth := NewAuth(fakeSessionReader{})
	handler := auth.RequireAuth(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler must not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/logout", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRequireAuthSessionNotFoundReturnsUnauthorized(t *testing.T) {
	auth := NewAuth(fakeSessionReader{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{}, session.ErrSessionNotFound
		},
	})
	handler := auth.RequireAuth(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler must not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "session-123"})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRequireAuthWithValidSessionCallsNext(t *testing.T) {
	auth := NewAuth(fakeSessionReader{
		getFn: func(_ context.Context, sessionID string) (session.SessionData, error) {
			if sessionID != "session-123" {
				t.Fatalf("unexpected session id: %s", sessionID)
			}
			return session.SessionData{
				UserID:    "user-1",
				Role:      "customer",
				CreatedAt: time.Now().UTC(),
			}, nil
		},
	})

	nextCalled := false
	handler := auth.RequireAuth(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		nextCalled = true
		sessionID, ok := SessionIDFromContext(r.Context())
		if !ok || sessionID != "session-123" {
			t.Fatalf("missing session id in context: %v %s", ok, sessionID)
		}
		data, ok := SessionDataFromContext(r.Context())
		if !ok || data.UserID != "user-1" {
			t.Fatalf("missing session data in context: %v %#v", ok, data)
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "session-123"})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}
}

func TestRequireAuthGetErrorReturnsUnauthorized(t *testing.T) {
	auth := NewAuth(fakeSessionReader{
		getFn: func(_ context.Context, _ string) (session.SessionData, error) {
			return session.SessionData{}, errors.New("redis down")
		},
	})
	handler := auth.RequireAuth(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler must not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "session-123"})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}
