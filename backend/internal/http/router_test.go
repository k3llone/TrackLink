package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tracklink/internal/config"
	"tracklink/internal/platform/session"
)

type fakeRouterSessionStore struct {
	sessions map[string]session.SessionData
}

func (f *fakeRouterSessionStore) Create(_ context.Context, sessionID string, data session.SessionData, _ time.Duration) error {
	if f.sessions == nil {
		f.sessions = map[string]session.SessionData{}
	}
	f.sessions[sessionID] = data
	return nil
}

func (f *fakeRouterSessionStore) Get(_ context.Context, sessionID string) (session.SessionData, error) {
	if data, ok := f.sessions[sessionID]; ok {
		return data, nil
	}
	return session.SessionData{}, session.ErrSessionNotFound
}

func (f *fakeRouterSessionStore) Delete(_ context.Context, sessionID string) error {
	delete(f.sessions, sessionID)
	return nil
}

func TestLogoutRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestLogoutRouteRepeatLogoutReturnsUnauthorized(t *testing.T) {
	store := &fakeRouterSessionStore{
		sessions: map[string]session.SessionData{
			"session-123": {
				UserID: "user-1",
				Role:   "customer",
			},
		},
	}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
		},
	})

	firstReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	firstReq.AddCookie(&http.Cookie{Name: "tracklink_session", Value: "session-123"})
	firstResp := httptest.NewRecorder()
	router.ServeHTTP(firstResp, firstReq)

	if firstResp.Code != http.StatusNoContent {
		t.Fatalf("expected first logout status %d, got %d", http.StatusNoContent, firstResp.Code)
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	secondReq.AddCookie(&http.Cookie{Name: "tracklink_session", Value: "session-123"})
	secondResp := httptest.NewRecorder()
	router.ServeHTTP(secondResp, secondReq)

	if secondResp.Code != http.StatusUnauthorized {
		t.Fatalf("expected second logout status %d, got %d", http.StatusUnauthorized, secondResp.Code)
	}
}
