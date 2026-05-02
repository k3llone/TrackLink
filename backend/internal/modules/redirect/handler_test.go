package redirect

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestHandlerRedirectByCodeResolvesActiveLink(t *testing.T) {
	clicked := false
	touched := false
	repo := fakeRepository{
		findByCodeOrAliasFn: func(_ context.Context, code string) (Link, error) {
			if code != "promo" {
				t.Fatalf("expected code promo, got %s", code)
			}
			return Link{
				ID:        "link-1",
				Code:      "abc123",
				TargetURL: "https://example.com/promo",
				Status:    StatusActive,
			}, nil
		},
		createClickEventFn: func(_ context.Context, event ClickEvent) error {
			clicked = true
			if event.LinkID != "link-1" {
				t.Fatalf("expected link id link-1, got %s", event.LinkID)
			}
			return nil
		},
		touchActiveLinkFn: func(_ context.Context, linkID string, _ time.Time) error {
			touched = true
			if linkID != "link-1" {
				t.Fatalf("expected link id link-1, got %s", linkID)
			}
			return nil
		},
	}
	handler := NewHandler(NewService(repo))
	router := chi.NewRouter()
	router.Get("/{code}", handler.RedirectByCode)

	req := httptest.NewRequest(http.MethodGet, "/promo", nil)
	req.Header.Set("Referer", "https://source.example.com")
	req.Header.Set("User-Agent", "test-agent")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected status %d, got %d", http.StatusFound, rr.Code)
	}
	if location := rr.Header().Get("Location"); location != "https://example.com/promo" {
		t.Fatalf("expected location https://example.com/promo, got %s", location)
	}
	if !clicked {
		t.Fatal("expected click event to be created")
	}
	if !touched {
		t.Fatal("expected active link click stats to be updated")
	}
}
