package admin

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"tracklink/internal/shared"
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

func (h *Handler) ListLinks(w http.ResponseWriter, r *http.Request) {
	adminUserID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(adminUserID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	query, fields := parseListLinksQuery(r)
	if len(fields) > 0 {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid query parameters", fields)
		return
	}

	items, pagination, serviceFields, err := h.service.List(r.Context(), adminUserID, query)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid query parameters", serviceFields)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	respItems := make([]AdminLink, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, mapAdminLink(item, h.publicURL))
	}

	writeJSON(w, http.StatusOK, AdminLinkListResponse{
		Items:      respItems,
		Pagination: mapPagination(pagination),
	})
}

func (h *Handler) BlockLink(w http.ResponseWriter, r *http.Request) {
	adminUserID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(adminUserID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	linkID := strings.TrimSpace(chi.URLParam(r, "linkId"))
	var req AdminBlockLinkRequest
	if err := decodeOptionalJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", nil)
		return
	}

	link, fields, err := h.service.Block(r.Context(), adminUserID, linkID, req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", fields)
			return
		}
		if errors.Is(err, ErrLinkNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Resource not found", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	writeJSON(w, http.StatusOK, mapAdminLink(link, h.publicURL))
}

func (h *Handler) DeactivateLink(w http.ResponseWriter, r *http.Request) {
	adminUserID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(adminUserID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	linkID := strings.TrimSpace(chi.URLParam(r, "linkId"))
	link, fields, err := h.service.Deactivate(r.Context(), adminUserID, linkID)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", fields)
			return
		}
		if errors.Is(err, ErrLinkNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Resource not found", nil)
			return
		}
		if errors.Is(err, ErrStatusChangeNotAllowed) {
			writeError(w, http.StatusConflict, "status_change_not_allowed", "Status change is not allowed", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	writeJSON(w, http.StatusOK, mapAdminLink(link, h.publicURL))
}

func decodeOptionalJSONBody(r *http.Request, dst any) error {
	if r.Body == nil {
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	return nil
}

func parseListLinksQuery(r *http.Request) (ListLinksQuery, map[string]string) {
	q := r.URL.Query()
	page := defaultListPage
	pageSize := defaultListPageSize
	fields := map[string]string{}

	if rawPage := strings.TrimSpace(q.Get("page")); rawPage != "" {
		v, err := strconv.Atoi(rawPage)
		if err != nil {
			fields["page"] = "Page must be an integer"
		} else if v < 1 {
			fields["page"] = "Page must be greater than or equal to 1"
		} else {
			page = v
		}
	}

	if rawPageSize := strings.TrimSpace(q.Get("pageSize")); rawPageSize != "" {
		v, err := strconv.Atoi(rawPageSize)
		if err != nil {
			fields["pageSize"] = "Page size must be an integer"
		} else if v < 1 {
			fields["pageSize"] = "Page size must be greater than or equal to 1"
		} else if v > maxListPageSize {
			pageSize = maxListPageSize
		} else {
			pageSize = v
		}
	}

	return ListLinksQuery{
		Page:     page,
		PageSize: pageSize,
		Q:        strings.TrimSpace(q.Get("q")),
	}, fields
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
