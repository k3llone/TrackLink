package shared

import (
	"context"

	"tracklink/internal/platform/session"
)

type contextKey string

const (
	sessionIDContextKey   contextKey = "session_id"
	sessionDataContextKey contextKey = "session_data"
)

func WithCurrentSession(ctx context.Context, sessionID string, data session.SessionData) context.Context {
	ctxWithSession := context.WithValue(ctx, sessionIDContextKey, sessionID)
	return context.WithValue(ctxWithSession, sessionDataContextKey, data)
}

func SessionIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(sessionIDContextKey).(string)
	return v, ok
}

func SessionDataFromContext(ctx context.Context) (session.SessionData, bool) {
	v, ok := ctx.Value(sessionDataContextKey).(session.SessionData)
	return v, ok
}

func CurrentUserFromContext(ctx context.Context) (string, string, bool) {
	data, ok := SessionDataFromContext(ctx)
	if !ok {
		return "", "", false
	}

	return data.UserID, data.Role, true
}
