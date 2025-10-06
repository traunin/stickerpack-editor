package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/queue"
)

const (
	baseRoute        = "/api"
	publicPacksRoute = "/public/packs"
	userPacksRoute   = "/user/packs"
	userPackRoute    = "/user/packs/"
	sessionRoute     = "/session"
	thumbnailRoute   = "/thumbnail"
	jobStatusRoute   = "/job/"
	queueStatusRoute = "/queue"
)

var noAuthRoutes = []NoAuthRoute{
	{Path: baseRoute + sessionRoute, Method: http.MethodPost, PrefixMatch: false},
	{Path: baseRoute + publicPacksRoute, Method: http.MethodGet, PrefixMatch: false},
	{Path: baseRoute + userPackRoute, Method: http.MethodHead, PrefixMatch: true},
	{Path: baseRoute + thumbnailRoute, Method: http.MethodGet, PrefixMatch: false},
	// {Path: baseRoute + jobStatusRoute, Method: http.MethodGet, PrefixMatch: true},
	{Path: baseRoute + queueStatusRoute, Method: http.MethodGet, PrefixMatch: false},
	{Path: "", Method: http.MethodOptions, PrefixMatch: true}, // preflight
}

var jobQueue *queue.Queue

func withCORS(h http.Handler) http.Handler {
	domain := config.Load().Domain()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domain)
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS, HEAD",
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

		h.ServeHTTP(w, r)
	})
}

func withContentTypeJSON(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, thumbnailRoute) &&
			!strings.HasPrefix(r.URL.Path, baseRoute+jobStatusRoute) {
			w.Header().Set("Content-Type", "application/json")
		}
		h.ServeHTTP(w, r)
	})
}

func SetupRouter() http.Handler {
	cfg := config.Load()
	jobQueue = queue.NewQueue(cfg.QueueWorkers())

	mux := http.NewServeMux()
	api := http.NewServeMux()

	api.HandleFunc(sessionRoute, sessionHandler)
	api.HandleFunc(publicPacksRoute, publicPacksHandler)
	api.HandleFunc(userPackRoute, userPackHandler)
	api.HandleFunc(userPacksRoute, userPacksHandler)
	api.HandleFunc(thumbnailRoute, thumbnailHandler)
	api.HandleFunc(jobStatusRoute, jobStatusHandler)
	api.HandleFunc(queueStatusRoute, queueStatsHandler)

	mux.Handle(baseRoute+"/", http.StripPrefix(baseRoute, api))

	auth := JWTMiddleware(noAuthRoutes, mux)
	contentType := withContentTypeJSON(auth)
	cors := withCORS(contentType)

	return cors
}

func publicPacksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPublicPacksHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userPackHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, userPackRoute)
	if name == "" {
		http.Error(w, "Missing pack name", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		deletePackHandler(w, r, name)
	case http.MethodHead:
		nameExistsHandler(w, r, name)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userPacksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUserPacksHandler(w, r)
	case http.MethodPost:
		createPackHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodPost:
		createSessionHandler(w, r)
	case http.MethodDelete:
		deleteSessionHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobID := strings.TrimPrefix(r.URL.Path, jobStatusRoute)
	if jobID == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	jobQueue.SSEHandler(w, r, jobID)
}

func queueStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := jobQueue.GetQueueStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
