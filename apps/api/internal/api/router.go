package api

import (
	"net/http"
	"strings"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

const (
	packsRoute = "/packs"
	packRoute  = "/packs/"
)

func withCORS(h http.Handler) http.Handler {
	domain := config.Load().DomainCORS
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domain)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	api := http.NewServeMux()
	api.HandleFunc(packsRoute, packsHandler)
	api.HandleFunc(packRoute, packHandler)

	mux.Handle("/api/", http.StripPrefix("/api", api))

	return withCORS(mux)
}

func packsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createPackHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func packHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, packRoute)
	if id == "" {
		http.Error(w, "Missing pack ID", http.StatusBadRequest)
		return
	}
	
	// TODO delete packs with id
	switch r.Method {
	case http.MethodDelete:
		deletePackHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
