package httpserver

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"tracklink/internal/config"
	"tracklink/internal/modules/accounts"
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

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestLogoutRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(testLogger(), Deps{
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
	router := NewRouter(testLogger(), Deps{
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
	router := NewRouter(testLogger(), Deps{
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
	router := NewRouter(testLogger(), Deps{
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
	router := NewRouter(testLogger(), Deps{
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
	router := NewRouter(testLogger(), Deps{
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
	router := NewRouter(testLogger(), Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/s/some-code", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusUnauthorized {
		t.Fatalf("expected redirect route to be public, got %d", rr.Code)
	}
}

func TestCreateLinkRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(testLogger(), Deps{
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
	router := NewRouter(testLogger(), Deps{
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

func TestUpdateLinkStatusRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(testLogger(), Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
			PublicURL:           "https://tracklink.example.com",
		},
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/links/7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a/status", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestDeleteLinkRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(testLogger(), Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
			PublicURL:           "https://tracklink.example.com",
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestDashboardRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(testLogger(), Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
			PublicURL:           "https://tracklink.example.com",
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestLinkAnalyticsRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(testLogger(), Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
			PublicURL:           "https://tracklink.example.com",
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links/7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a/analytics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRecentClicksRouteRequiresAuth(t *testing.T) {
	store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
	router := NewRouter(testLogger(), Deps{
		Sessions: store,
		Config: config.Config{
			SessionTTL:          24 * time.Hour,
			SessionCookieSecure: false,
			PublicURL:           "https://tracklink.example.com",
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/links/7607f3ca-90d7-4c47-b2f7-f968ad1f5f9a/clicks", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAdminRouteRequiresAdmin(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		sessionData session.SessionData
		wantStatus  int
	}{
		{
			name:       "unauthorized without session",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "forbidden for customer",
			sessionID: "session-customer",
			sessionData: session.SessionData{
				UserID: "user-1",
				Role:   accounts.RoleCustomer,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:      "allowed for admin",
			sessionID: "session-admin",
			sessionData: session.SessionData{
				UserID: "admin-1",
				Role:   accounts.RoleAdmin,
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
			if tt.sessionID != "" {
				store.sessions[tt.sessionID] = tt.sessionData
			}

			router := newRouter(testLogger(), Deps{
				Sessions: store,
				Config: config.Config{
					SessionTTL:          24 * time.Hour,
					SessionCookieSecure: false,
				},
			}, func(r chi.Router) {
				r.Get("/probe", func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/probe", nil)
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{Name: "tracklink_session", Value: tt.sessionID})
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}
		})
	}
}

func TestAdminEndpointsAreMountedAndProtected(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		sessionID   string
		sessionData session.SessionData
		wantStatus  int
	}{
		{
			name:       "list returns unauthorized without session",
			method:     http.MethodGet,
			path:       "/api/v1/admin/links",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "list returns forbidden for customer",
			method:    http.MethodGet,
			path:      "/api/v1/admin/links",
			sessionID: "session-customer",
			sessionData: session.SessionData{
				UserID: "user-1",
				Role:   accounts.RoleCustomer,
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "block returns unauthorized without session",
			method:     http.MethodPatch,
			path:       "/api/v1/admin/links/link-1/block",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "block returns forbidden for customer",
			method:    http.MethodPatch,
			path:      "/api/v1/admin/links/link-1/block",
			sessionID: "session-customer",
			sessionData: session.SessionData{
				UserID: "user-1",
				Role:   accounts.RoleCustomer,
			},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeRouterSessionStore{sessions: map[string]session.SessionData{}}
			if tt.sessionID != "" {
				store.sessions[tt.sessionID] = tt.sessionData
			}

			router := NewRouter(testLogger(), Deps{
				Sessions: store,
				Config: config.Config{
					SessionTTL:          24 * time.Hour,
					SessionCookieSecure: false,
					PublicURL:           "https://tracklink.example.com",
				},
			})

			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{Name: "tracklink_session", Value: tt.sessionID})
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}
		})
	}
}
