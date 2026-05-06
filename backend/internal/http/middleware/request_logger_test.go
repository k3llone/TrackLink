package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggerLogsMethodPathStatusAndDuration(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	mw := RequestLogger(log)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "ok")
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", strings.NewReader(`{"title":"x"}`))
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "session-secret"})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	out := buf.String()
	if !strings.Contains(out, "msg=http_request") {
		t.Fatalf("expected http_request log, got %q", out)
	}
	if !strings.Contains(out, "method=POST") {
		t.Fatalf("expected method in log, got %q", out)
	}
	if !strings.Contains(out, "path=/api/v1/links") {
		t.Fatalf("expected path in log, got %q", out)
	}
	if !strings.Contains(out, "status=201") {
		t.Fatalf("expected status in log, got %q", out)
	}
	if !strings.Contains(out, "duration=") {
		t.Fatalf("expected duration in log, got %q", out)
	}
	if strings.Contains(out, "Authorization") || strings.Contains(out, "test-token") {
		t.Fatalf("authorization header leaked into logs: %q", out)
	}
	if strings.Contains(out, "session-secret") {
		t.Fatalf("cookie value leaked into logs: %q", out)
	}
	if strings.Contains(out, "title") {
		t.Fatalf("request body leaked into logs: %q", out)
	}
}

func TestRequestLoggerDefaultsStatusToOK(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	mw := RequestLogger(log)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	out := buf.String()
	if !strings.Contains(out, "status=200") {
		t.Fatalf("expected default status=200 in log, got %q", out)
	}
}
