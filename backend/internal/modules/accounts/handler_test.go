package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"tracklink/internal/platform/session"
)

type fakeSessionStore struct {
	createFn func(ctx context.Context, sessionID string, data session.SessionData, ttl time.Duration) error
}

func (f fakeSessionStore) Create(ctx context.Context, sessionID string, data session.SessionData, ttl time.Duration) error {
	if f.createFn == nil {
		return nil
	}
	return f.createFn(ctx, sessionID, data, ttl)
}

func TestHandlerRegisterSuccess(t *testing.T) {
	repo := fakeRepo{
		createFn: func(_ context.Context, user *User) error {
			user.ID = "7c2f61cd-1067-4b49-85c2-1f0d85a60401"
			user.CreatedAt = time.Date(2026, 4, 30, 17, 10, 0, 0, time.UTC)
			return nil
		},
	}

	handler := NewHandler(NewService(repo), fakeSessionStore{}, CookieSettings{})

	body := `{"email":"new@example.com","password":"StrongPass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"email":"new@example.com"`) {
		t.Fatalf("expected email in response, got %s", rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "password_hash") {
		t.Fatalf("response must not expose password_hash: %s", rr.Body.String())
	}
}

func TestHandlerRegisterInvalidBody(t *testing.T) {
	repo := fakeRepo{
		createFn: func(_ context.Context, _ *User) error { return nil },
	}

	handler := NewHandler(NewService(repo), fakeSessionStore{}, CookieSettings{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{"email"`))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

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

func TestHandlerRegisterDuplicateEmail(t *testing.T) {
	repo := fakeRepo{
		createFn: func(_ context.Context, _ *User) error {
			return ErrEmailAlreadyExists
		},
	}

	handler := NewHandler(NewService(repo), fakeSessionStore{}, CookieSettings{})
	body := `{"email":"new@example.com","password":"StrongPass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "email_already_exists" {
		t.Fatalf("unexpected error code: %s", resp.Error.Code)
	}
}

func TestHandlerLoginSuccess(t *testing.T) {
	passwordHash, err := hashPassword("StrongPass123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := fakeRepo{
		findByEmailFn: func(_ context.Context, _ string) (User, error) {
			return User{
				ID:           "4a9550d6-c2df-4adb-8b44-f85e6f02177f",
				Email:        "new@example.com",
				PasswordHash: passwordHash,
				Role:         RoleCustomer,
				CreatedAt:    time.Date(2026, 4, 30, 17, 10, 0, 0, time.UTC),
			}, nil
		},
	}

	var savedSessionID string
	handler := NewHandler(NewService(repo), fakeSessionStore{
		createFn: func(_ context.Context, sessionID string, data session.SessionData, ttl time.Duration) error {
			savedSessionID = sessionID
			if data.UserID == "" {
				t.Fatal("user ID must be present in session data")
			}
			if ttl != 24*time.Hour {
				t.Fatalf("unexpected ttl: %v", ttl)
			}
			return nil
		},
	}, CookieSettings{
		Name:   "tracklink_session",
		TTL:    24 * time.Hour,
		Secure: false,
	})

	body := `{"email":"new@example.com","password":"StrongPass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if savedSessionID == "" {
		t.Fatal("session id must be generated")
	}
	if strings.Contains(rr.Body.String(), "password_hash") {
		t.Fatalf("response must not expose password_hash: %s", rr.Body.String())
	}

	cookies := rr.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected at least one cookie")
	}

	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "tracklink_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("tracklink_session cookie not found")
	}
	if sessionCookie.Value == "" {
		t.Fatal("session cookie must have value")
	}
}

func TestHandlerLoginUnauthorized(t *testing.T) {
	repo := fakeRepo{
		findByEmailFn: func(_ context.Context, _ string) (User, error) {
			return User{}, ErrUserNotFound
		},
	}

	handler := NewHandler(NewService(repo), fakeSessionStore{}, CookieSettings{})
	body := `{"email":"new@example.com","password":"WrongPass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "unauthorized" {
		t.Fatalf("unexpected error code: %s", resp.Error.Code)
	}
}

func TestHandlerLoginInvalidBody(t *testing.T) {
	handler := NewHandler(NewService(fakeRepo{}), fakeSessionStore{}, CookieSettings{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{"email"`))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandlerLoginSessionStoreError(t *testing.T) {
	passwordHash, err := hashPassword("StrongPass123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	repo := fakeRepo{
		findByEmailFn: func(_ context.Context, _ string) (User, error) {
			return User{
				ID:           "4a9550d6-c2df-4adb-8b44-f85e6f02177f",
				Email:        "new@example.com",
				PasswordHash: passwordHash,
				Role:         RoleCustomer,
			}, nil
		},
	}
	handler := NewHandler(NewService(repo), fakeSessionStore{
		createFn: func(_ context.Context, _ string, _ session.SessionData, _ time.Duration) error {
			return errors.New("redis unavailable")
		},
	}, CookieSettings{
		Name: "tracklink_session",
		TTL:  24 * time.Hour,
	})

	body := `{"email":"new@example.com","password":"StrongPass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
