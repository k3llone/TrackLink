package links

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	httpmiddleware "tracklink/internal/http/middleware"
)

type Handler struct {
	service   *Service
	publicURL string
}

func NewHandler(service *Service, publicURL string) *Handler {
	return &Handler{
		service:   service,
		publicURL: strings.TrimRight(strings.TrimSpace(publicURL), "/"),
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	sessionData, ok := httpmiddleware.SessionDataFromContext(r.Context())
	if !ok || strings.TrimSpace(sessionData.UserID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", nil)
		return
	}

	link, fields, err := h.service.Create(r.Context(), sessionData.UserID, req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", fields)
			return
		}

		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	writeJSON(w, http.StatusCreated, mapLinkToResponse(link, h.publicURL))
}

func mapLinkToResponse(link Link, publicURL string) LinkResponse {
	shortPath := link.Code
	if link.CustomAlias != nil && strings.TrimSpace(*link.CustomAlias) != "" {
		shortPath = strings.TrimSpace(*link.CustomAlias)
	}

	resp := LinkResponse{
		ID:          link.ID,
		OwnerID:     link.OwnerID,
		Code:        link.Code,
		CustomAlias: link.CustomAlias,
		ShortURL:    strings.TrimRight(publicURL, "/") + "/" + shortPath,
		TargetURL:   link.TargetURL,
		Status:      link.Status,
		TotalClicks: link.TotalClicks,
		CreatedAt:   link.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   link.UpdatedAt.UTC().Format(time.RFC3339),
	}

	if link.LastClickedAt != nil {
		ts := link.LastClickedAt.UTC().Format(time.RFC3339)
		resp.LastClickedAt = &ts
	}

	return resp
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
