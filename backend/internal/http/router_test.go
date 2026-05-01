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

func TestMeRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestMeRouteInvalidSessionReturnsUnauthorized(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.AddCookie(&http.Cookie{Name: "tracklink_session", Value: "missing-session"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestMeRouteValidSessionPassesMiddleware(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.AddCookie(&http.Cookie{Name: "tracklink_session", Value: "session-123"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusUnauthorized {
		t.Fatalf("expected request to pass middleware, got status %d", rr.Code)
	}
}

func TestPublicRoutesAreNotBlockedByAuthMiddleware(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
		},
	})

	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
	registerResp := httptest.NewRecorder()
	router.ServeHTTP(registerResp, registerReq)
	if registerResp.Code == http.StatusUnauthorized {
		t.Fatalf("expected public register route to be accessible, got %d", registerResp.Code)
	}

	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)
	if loginResp.Code == http.StatusUnauthorized {
		t.Fatalf("expected public login route to be accessible, got %d", loginResp.Code)
	}
}

func TestPublicRedirectRouteIsNotBlockedByAuthMiddleware(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/some-code", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusUnauthorized {
		t.Fatalf("expected redirect route to be public, got %d", rr.Code)
	}
}

func TestCreateLinkRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
			PublicURL:           "https://tracklink.example.com",
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestListLinksRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
			PublicURL:           "https://tracklink.example.com",
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}
