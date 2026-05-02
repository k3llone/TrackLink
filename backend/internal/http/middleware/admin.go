package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"tracklink/internal/modules/accounts"
	"tracklink/internal/shared"
)

func (a *Auth) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, role, ok := shared.CurrentUserFromContext(r.Context())
		if !ok || strings.TrimSpace(userID) == "" {
			writeUnauthorized(w)
			return
		}

		if strings.TrimSpace(role) != accounts.RoleAdmin {
			writeForbidden(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: errorBody{
			Code:    "forbidden",
			Message: "Forbidden",
		},
	})
}
