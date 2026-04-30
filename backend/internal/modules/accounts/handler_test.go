package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandlerRegisterSuccess(t *testing.T) {
	repo := fakeRepo{
		createFn: func(_ context.Context, user *User) error {
			user.ID = "7c2f61cd-1067-4b49-85c2-1f0d85a60401"
			user.CreatedAt = time.Date(2026, 4, 30, 17, 10, 0, 0, time.UTC)
			return nil
		},
	}

	handler := NewHandler(NewService(repo))

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

	handler := NewHandler(NewService(repo))
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
