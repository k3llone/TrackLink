package httpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"tracklink/internal/config"
	httpmiddleware "tracklink/internal/http/middleware"
	"tracklink/internal/modules/accounts"
	"tracklink/internal/modules/admin"
	"tracklink/internal/modules/analytics"
	"tracklink/internal/modules/links"
	"tracklink/internal/modules/redirect"
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

func NewRouter(log *slog.Logger, deps Deps) *chi.Mux {
	return newRouter(log, deps, nil)
}

func newRouter(log *slog.Logger, deps Deps, registerAdminRoutes func(chi.Router)) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(httpmiddleware.RequestLogger(log))

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
	adminRepo := admin.NewGormRepository(deps.DB)
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService, deps.Config.PublicURL)
	redirectRepo := redirect.NewGormRepository(deps.DB)
	analyticsRepo := analytics.NewGormRepository(deps.DB)
	redirectService := redirect.NewService(redirectRepo, analyticsRepo)
	redirectHandler := redirect.NewHandler(redirectService)
	analyticsService := analytics.NewService(analyticsRepo, deps.Config.PublicURL)
	analyticsHandler := analytics.NewHandler(analyticsService)
	adminV1 := chi.NewRouter()
	adminV1.Use(authMiddleware.RequireAuth)
	adminV1.Use(authMiddleware.RequireAdmin)
	adminV1.Get("/links", adminHandler.ListLinks)
	adminV1.Patch("/links/{linkId}/block", adminHandler.BlockLink)
	adminV1.Patch("/links/{linkId}/deactivate", adminHandler.DeactivateLink)
	if registerAdminRoutes != nil {
		registerAdminRoutes(adminV1)
	}
	apiV1.Post("/auth/register", accountHandler.Register)
	apiV1.Post("/auth/login", accountHandler.Login)
	apiV1.With(authMiddleware.RequireAuth).Post("/auth/logout", accountHandler.Logout)
	apiV1.With(authMiddleware.RequireAuth).Get("/me", accountHandler.Me)
	apiV1.With(authMiddleware.RequireAuth).Get("/links", linkHandler.List)
	apiV1.With(authMiddleware.RequireAuth).Post("/links", linkHandler.Create)
	apiV1.With(authMiddleware.RequireAuth).Patch("/links/{linkId}/status", linkHandler.UpdateStatus)
	apiV1.With(authMiddleware.RequireAuth).Delete("/links/{linkId}", linkHandler.Delete)
	apiV1.With(authMiddleware.RequireAuth).Get("/dashboard", analyticsHandler.Dashboard)
	apiV1.With(authMiddleware.RequireAuth).Get("/links/{linkId}/analytics", analyticsHandler.LinkAnalytics)
	apiV1.With(authMiddleware.RequireAuth).Get("/links/{linkId}/clicks", analyticsHandler.RecentClicks)
	apiV1.Mount("/admin", adminV1)
	r.Mount("/api/v1", apiV1)

	r.Get("/s/{code}", redirectHandler.RedirectByCode)

	return r
}

func defaultSessionTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return 24 * time.Hour
	}

	return ttl
}
