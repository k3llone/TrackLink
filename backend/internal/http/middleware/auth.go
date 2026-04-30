package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"tracklink/internal/platform/session"
)

const SessionCookieName = "tracklink_session"

type SessionReader interface {
	Get(ctx context.Context, sessionID string) (session.SessionData, error)
}

type Auth struct {
	sessions SessionReader
}

type contextKey string

const (
	sessionIDContextKey   contextKey = "session_id"
	sessionDataContextKey contextKey = "session_data"
)

func NewAuth(sessions SessionReader) *Auth {
	return &Auth{sessions: sessions}
}

func (a *Auth) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.sessions == nil {
			writeUnauthorized(w)
			return
		}

		cookie, err := r.Cookie(SessionCookieName)
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			writeUnauthorized(w)
			return
		}

		sessionID := strings.TrimSpace(cookie.Value)
		data, err := a.sessions.Get(r.Context(), sessionID)
		if err != nil {
			writeUnauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), sessionIDContextKey, sessionID)
		ctx = context.WithValue(ctx, sessionDataContextKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SessionIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(sessionIDContextKey).(string)
	return v, ok
}

func SessionDataFromContext(ctx context.Context) (session.SessionData, bool) {
	v, ok := ctx.Value(sessionDataContextKey).(session.SessionData)
	return v, ok
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: errorBody{
			Code:    "unauthorized",
			Message: "Unauthorized",
		},
	})
}
