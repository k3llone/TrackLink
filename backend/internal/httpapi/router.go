package httpapi

import (
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
	_ = deps

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	return r
}
