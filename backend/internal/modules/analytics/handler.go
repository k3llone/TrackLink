package analytics

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) LinkAnalytics(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	query, fields := parseLinkAnalyticsQuery(r)
	if len(fields) > 0 {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid query parameters", fields)
		return
	}

	linkID := strings.TrimSpace(chi.URLParam(r, "linkId"))
	response, serviceFields, err := h.service.LoadLinkAnalytics(r.Context(), userID, linkID, query)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid query parameters", serviceFields)
			return
		}
		if errors.Is(err, ErrLinkNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Resource not found", nil)
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

func parseLinkAnalyticsQuery(r *http.Request) (LinkAnalyticsQuery, map[string]string) {
	q := r.URL.Query()
	fields := map[string]string{}
	groupBy := strings.TrimSpace(q.Get("groupBy"))
	if groupBy == "" {
		groupBy = defaultAnalyticsGroupBy
	}
	if groupBy != GroupByHour && groupBy != GroupByDay {
		fields["groupBy"] = "Group by must be one of: hour, day"
	}

	return LinkAnalyticsQuery{
		GroupBy: groupBy,
	}, fields
}
