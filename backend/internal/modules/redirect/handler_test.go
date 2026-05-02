package redirect

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHandlerRedirectByCodeReturnsNotFound(t *testing.T) {
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, _ string) (Link, error) {
			return Link{}, ErrLinkNotFound
		},
	}
	handler := NewHandler(NewService(repo))
	router := chi.NewRouter()
	router.Get("/{code}", handler.RedirectByCode)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Link not found") {
		t.Fatalf("expected html not found message, got %s", rr.Body.String())
	}
}

func TestHandlerRedirectByCodeResolvesTarget(t *testing.T) {
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != "promo" {
				t.Fatalf("expected code promo, got %s", code)
			}
			return Link{
				ID:        "link-1",
				Code:      "abc123",
				TargetURL: "https://example.com/promo",
				Status:    "active",
			}, nil
		},
	}
	handler := NewHandler(NewService(repo))
	router := chi.NewRouter()
	router.Get("/{code}", handler.RedirectByCode)

	req := httptest.NewRequest(http.MethodGet, "/promo", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "https://example.com/promo") {
		t.Fatalf("expected resolved target url in body, got %s", rr.Body.String())
	}
}
