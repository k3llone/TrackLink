package redirect

import (
	"net/http"
	"time"

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
	result, err := h.service.ResolveAndTrack(r.Context(), code, RequestMeta{
		Referrer:  r.Referer(),
		UserAgent: r.UserAgent(),
		ClickedAt: time.Now().UTC(),
	})
	if err != nil {
		writeHTML(w, http.StatusInternalServerError, "Internal error", "Unable to process redirect request.")
		return
	}

	switch result.Kind {
	case ResultKindNotFound:
		writeHTML(w, http.StatusNotFound, "Link not found", "The requested short link does not exist.")
	case ResultKindRedirect:
		http.Redirect(w, r, result.TargetURL, http.StatusFound)
	default:
		writeHTML(w, http.StatusConflict, "Link unavailable", "The requested link is temporarily unavailable.")
	}
}

func writeHTML(w http.ResponseWriter, status int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte("<!doctype html><html><head><meta charset=\"utf-8\"><title>" + title + "</title></head><body><h1>" + title + "</h1><p>" + message + "</p></body></html>"))
}
