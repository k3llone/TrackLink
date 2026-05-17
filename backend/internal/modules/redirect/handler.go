package redirect

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

const (
	linkUnavailableNotFoundPath = "/link-unavailable/not-found"
	linkUnavailableBlockedPath  = "/link-unavailable/blocked"
	linkUnavailableInactivePath = "/link-unavailable/inactive"
	linkUnavailableDeletedPath  = "/link-unavailable/deleted"
	linkUnavailableGonePath     = "/link-unavailable/gone"
)

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
		redirectToUnavailablePage(w, r, linkUnavailableNotFoundPath)
	case ResultKindRedirect:
		http.Redirect(w, r, result.TargetURL, http.StatusFound)
	default:
		handleUnavailable(w, r, result)
	}
}

func handleUnavailable(w http.ResponseWriter, r *http.Request, result ResolveResult) {
	switch {
	case result.Status == StatusBlocked:
		redirectToUnavailablePage(w, r, linkUnavailableBlockedPath)
	case result.Status == StatusInactive:
		redirectToUnavailablePage(w, r, linkUnavailableInactivePath)
	case result.Status == StatusDeleted || result.Deleted:
		redirectToUnavailablePage(w, r, linkUnavailableDeletedPath)
	default:
		redirectToUnavailablePage(w, r, linkUnavailableGonePath)
	}
}

func redirectToUnavailablePage(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
}

func writeHTML(w http.ResponseWriter, status int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte("<!doctype html><html><head><meta charset=\"utf-8\"><title>" + title + "</title></head><body><h1>" + title + "</h1><p>" + message + "</p></body></html>"))
}
