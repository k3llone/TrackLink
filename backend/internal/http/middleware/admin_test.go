package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tracklink/internal/modules/accounts"
	"tracklink/internal/platform/session"
	"tracklink/internal/shared"
)

func TestRequireAdminWithoutCurrentUserReturnsUnauthorized(t *testing.T) {
	auth := NewAuth(nil)
	handler := auth.RequireAdmin(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler must not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/probe", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "unauthorized" || resp.Error.Message != "Unauthorized" {
		t.Fatalf("unexpected error body: %#v", resp.Error)
	}
}

func TestRequireAdminWithCustomerRoleReturnsForbidden(t *testing.T) {
	auth := NewAuth(nil)
	handler := auth.RequireAdmin(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler must not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/probe", nil)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-123", session.SessionData{
		UserID:    "user-1",
		Role:      accounts.RoleCustomer,
		CreatedAt: time.Now().UTC(),
	}))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "forbidden" || resp.Error.Message != "Forbidden" {
		t.Fatalf("unexpected error body: %#v", resp.Error)
	}
}

func TestRequireAdminWithAdminRoleCallsNext(t *testing.T) {
	auth := NewAuth(nil)
	nextCalled := false
	handler := auth.RequireAdmin(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		nextCalled = true

		userID, role, ok := shared.CurrentUserFromContext(r.Context())
		if !ok || userID != "admin-1" || role != accounts.RoleAdmin {
			t.Fatalf("unexpected user in context: ok=%v userID=%q role=%q", ok, userID, role)
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/probe", nil)
	req = req.WithContext(shared.WithCurrentSession(req.Context(), "session-123", session.SessionData{
		UserID:    "admin-1",
		Role:      accounts.RoleAdmin,
		CreatedAt: time.Now().UTC(),
	}))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}
}
