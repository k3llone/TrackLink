package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"tracklink/internal/platform/session"
)

type Deps struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions *session.RedisStore
}

func NewRouter(deps Deps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	apiV1 := chi.NewRouter()
	apiV1.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
	r.Mount("/api/v1", apiV1)

	r.Get("/{code}", func(w http.ResponseWriter, _ *http.Request) {
		_ = deps
		http.Error(w, "redirect handler is not implemented yet", http.StatusNotImplemented)
	})

	return r
}
