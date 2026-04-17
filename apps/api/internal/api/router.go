package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/db"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/queue"
)

const (
	baseRoute        = "/api"
	publicPacksRoute = "/public/packs"
	userPacksRoute   = "/user/packs"
	userPackRoute    = "/user/packs/"
	sessionRoute     = "/session"
	mediaRoute       = "/media"
	jobStatusRoute   = "/job/"
	queueStatusRoute = "/queue"
)

var noAuthRoutes = []NoAuthRoute{
	{Path: baseRoute + sessionRoute, Method: http.MethodPost, PrefixMatch: false},
	{Path: baseRoute + publicPacksRoute, Method: http.MethodGet, PrefixMatch: false},
	{Path: baseRoute + userPackRoute, Method: http.MethodHead, PrefixMatch: true},
	{Path: baseRoute + mediaRoute, Method: http.MethodGet, PrefixMatch: false},
	// {Path: baseRoute + jobStatusRoute, Method: http.MethodGet, PrefixMatch: true},
	{Path: baseRoute + queueStatusRoute, Method: http.MethodGet, PrefixMatch: false},
	{Path: "", Method: http.MethodOptions, PrefixMatch: true}, // preflight
}

type Handler struct {
	cfg   *config.Config
	db    *db.Postgres
	queue *queue.Queue
}

func withCORS(domain string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domain)
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD",
		)
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func withContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, mediaRoute) &&
			!strings.HasPrefix(r.URL.Path, baseRoute+jobStatusRoute) {
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

func SetupHandler(cfg *config.Config, dbConn *db.Postgres) http.Handler {
	h := &Handler{
		cfg:   cfg,
		db:    dbConn,
		queue: queue.NewQueue(cfg.QueueWorkers()),
	}

	mux := http.NewServeMux()
	api := http.NewServeMux()

	api.HandleFunc(sessionRoute, h.sessionHandler)
	api.HandleFunc(publicPacksRoute, h.publicPacksHandler)
	api.HandleFunc(userPackRoute, h.userPackHandler)
	api.HandleFunc(userPacksRoute, h.userPacksHandler)
	api.HandleFunc(mediaRoute, h.mediaHandler)
	api.HandleFunc(jobStatusRoute, h.jobStatusHandler)
	api.HandleFunc(queueStatusRoute, h.queueStatsHandler)

	mux.Handle(baseRoute+"/", http.StripPrefix(baseRoute, api))

	auth := h.jwtMiddleware(noAuthRoutes, mux)
	contentType := withContentTypeJSON(auth)
	cors := withCORS(cfg.Domain(), contentType)

	return cors
}

func (h *Handler) publicPacksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getPublicPacksHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) userPackHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, userPackRoute)
	if name == "" {
		http.Error(w, "Missing pack name", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.deletePackHandler(w, r, name)
	case http.MethodGet:
		h.getPackHandler(w, r, name)
	case http.MethodPatch:
		h.editPackHandler(w, r, name)
	case http.MethodHead:
		h.nameExistsHandler(w, r, name)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) userPacksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUserPacksHandler(w, r)
	case http.MethodPost:
		h.createPackHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) sessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		h.createSessionHandler(w, r)
	case http.MethodDelete:
		h.deleteSessionHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobID := strings.TrimPrefix(r.URL.Path, jobStatusRoute)
	if jobID == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	h.queue.SSEHandler(w, r, jobID)
}

func (h *Handler) queueStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := h.queue.GetQueueStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
