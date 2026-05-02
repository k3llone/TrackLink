package redirect

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RedirectByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	link, err := h.service.Resolve(r.Context(), code)
	if err != nil {
		if errors.Is(err, ErrLinkNotFound) {
			writeHTML(w, http.StatusNotFound, "Link not found", "The requested short link does not exist.")
			return
		}
		writeHTML(w, http.StatusInternalServerError, "Internal error", "Unable to process redirect request.")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("resolved target: %s", link.TargetURL)))
}

func writeHTML(w http.ResponseWriter, status int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte("<!doctype html><html><head><meta charset=\"utf-8\"><title>" + title + "</title></head><body><h1>" + title + "</h1><p>" + message + "</p></body></html>"))
}
