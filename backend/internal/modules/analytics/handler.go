package analytics

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"tracklink/internal/shared"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	response, fields, err := h.service.LoadDashboard(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid request", fields)
			return
		}

		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

type errorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

func writeError(w http.ResponseWriter, status int, code, message string, fields map[string]string) {
	writeJSON(w, status, errorResponse{
		Error: errorBody{
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
