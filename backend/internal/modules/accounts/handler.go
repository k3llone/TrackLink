package accounts

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"tracklink/internal/platform/session"
	"tracklink/internal/shared"
)

const defaultSessionCookieName = "tracklink_session"

type SessionStore interface {
	Create(ctx context.Context, sessionID string, data session.SessionData, ttl time.Duration) error
	Delete(ctx context.Context, sessionID string) error
}

type CookieSettings struct {
	Name   string
	TTL    time.Duration
	Secure bool
}

type Handler struct {
	service      *Service
	sessions     SessionStore
	cookieConfig CookieSettings
}

func NewHandler(service *Service, sessions SessionStore, cookieConfig CookieSettings) *Handler {
	if cookieConfig.Name == "" {
		cookieConfig.Name = defaultSessionCookieName
	}
	if cookieConfig.TTL <= 0 {
		cookieConfig.TTL = 24 * time.Hour
	}

	return &Handler{
		service:      service,
		sessions:     sessions,
		cookieConfig: cookieConfig,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", nil)
		return
	}

	user, fields, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", fields)
			return
		}
		if errors.Is(err, ErrConflict) {
			writeError(w, http.StatusConflict, "email_already_exists", "Email already exists", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	resp := RegisterResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		},
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", nil)
		return
	}

	user, fields, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", fields)
			return
		}
		if errors.Is(err, ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Invalid email or password", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	if h.sessions == nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	sessionData := session.SessionData{
		UserID:    user.ID,
		Role:      user.Role,
		CreatedAt: time.Now().UTC(),
	}
	if err := h.sessions.Create(r.Context(), sessionID, sessionData, h.cookieConfig.TTL); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieConfig.Name,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(h.cookieConfig.TTL.Seconds()),
		Expires:  time.Now().UTC().Add(h.cookieConfig.TTL),
		HttpOnly: true,
		Secure:   h.cookieConfig.Secure,
		SameSite: http.SameSiteLaxMode,
	})

	resp := AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		},
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if h.sessions == nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	sessionID, ok := shared.SessionIDFromContext(r.Context())
	if !ok || strings.TrimSpace(sessionID) == "" {
		cookie, err := r.Cookie(h.cookieConfig.Name)
		if err != nil || strings.TrimSpace(cookie.Value) == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
			return
		}
		sessionID = strings.TrimSpace(cookie.Value)
	}

	if err := h.sessions.Delete(r.Context(), sessionID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieConfig.Name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
		HttpOnly: true,
		Secure:   h.cookieConfig.Secure,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	user, err := h.service.CurrentUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	resp := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeError(w http.ResponseWriter, status int, code, message string, fields map[string]string) {
	writeJSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func generateSessionID() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	return hex.EncodeToString(raw), nil
}
