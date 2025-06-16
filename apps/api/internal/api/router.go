package api

import (
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
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
	api.HandleFunc("/create-pack", createPackHandler)
	api.HandleFunc("/delete-pack", deletePackHandler)

	mux.Handle("/api/", http.StripPrefix("/api", api))

	return withCORS(mux)
}
