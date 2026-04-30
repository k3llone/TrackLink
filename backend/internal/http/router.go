package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"tracklink/internal/config"
	"tracklink/internal/modules/accounts"
	"tracklink/internal/platform/session"
)

type Deps struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions *session.RedisStore
	Config   config.Config
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

	userRepo := accounts.NewGormUserRepository(deps.DB)
	accountService := accounts.NewService(userRepo)
	accountHandler := accounts.NewHandler(accountService, deps.Sessions, accounts.CookieSettings{
		Name:   "tracklink_session",
		TTL:    defaultSessionTTL(deps.Config.SessionTTL),
		Secure: deps.Config.SessionCookieSecure,
	})
	apiV1.Post("/auth/register", accountHandler.Register)
	apiV1.Post("/auth/login", accountHandler.Login)
	r.Mount("/api/v1", apiV1)

	r.Get("/{code}", func(w http.ResponseWriter, _ *http.Request) {
		_ = deps
		http.Error(w, "redirect handler is not implemented yet", http.StatusNotImplemented)
	})

	return r
}

func defaultSessionTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return 24 * time.Hour
	}

	return ttl
}
