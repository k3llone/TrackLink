package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
