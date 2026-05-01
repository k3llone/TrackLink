package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"tracklink/internal/config"
	httpmiddleware "tracklink/internal/http/middleware"
	"tracklink/internal/modules/accounts"
	"tracklink/internal/modules/links"
	"tracklink/internal/platform/session"
)

type SessionStore interface {
	Create(ctx context.Context, sessionID string, data session.SessionData, ttl time.Duration) error
	Get(ctx context.Context, sessionID string) (session.SessionData, error)
	Delete(ctx context.Context, sessionID string) error
}

type Deps struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Sessions SessionStore
	Config   config.Config
}

func NewRouter(deps Deps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Logger)

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
	authMiddleware := httpmiddleware.NewAuth(deps.Sessions)
	linkRepo := links.NewGormRepository(deps.DB)
	linkService := links.NewService(linkRepo)
	linkHandler := links.NewHandler(linkService, deps.Config.PublicURL)
	apiV1.Post("/auth/register", accountHandler.Register)
	apiV1.Post("/auth/login", accountHandler.Login)
	apiV1.With(authMiddleware.RequireAuth).Post("/auth/logout", accountHandler.Logout)
	apiV1.With(authMiddleware.RequireAuth).Get("/me", accountHandler.Me)
	apiV1.With(authMiddleware.RequireAuth).Get("/links", linkHandler.List)
	apiV1.With(authMiddleware.RequireAuth).Post("/links", linkHandler.Create)
	apiV1.With(authMiddleware.RequireAuth).Patch("/links/{linkId}/status", linkHandler.UpdateStatus)
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
