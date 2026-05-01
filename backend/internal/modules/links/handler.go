package links

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", nil)
		return
	}

	link, fields, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid request body", fields)
			return
		}
		if errors.Is(err, ErrAliasAlreadyExists) {
			writeError(w, http.StatusConflict, "custom_alias_already_exists", "Custom alias is already taken", nil)
			return
		}

		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	writeJSON(w, http.StatusCreated, mapLinkToResponse(link, h.publicURL))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := shared.CurrentUserFromContext(r.Context())
	if !ok || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", nil)
		return
	}

	query, fields := parseListLinksQuery(r)
	if len(fields) > 0 {
		writeError(w, http.StatusBadRequest, "validation_error", "Invalid query parameters", fields)
		return
	}

	links, pagination, serviceFields, err := h.service.List(r.Context(), userID, query)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			writeError(w, http.StatusBadRequest, "validation_error", "Invalid query parameters", serviceFields)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}

	items := make([]LinkResponse, 0, len(links))
	for _, link := range links {
		items = append(items, mapLinkToResponse(link, h.publicURL))
	}

	writeJSON(w, http.StatusOK, LinkListResponse{
		Items:      items,
		Pagination: pagination,
	})
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
		Status:   strings.TrimSpace(q.Get("status")),
	}, fields
}
