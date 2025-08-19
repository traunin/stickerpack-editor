package api

import (
	"net/http"
	"strings"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
)

const (
	baseRoute        = "/api"
	publicPacksRoute = "/public/packs"
	userPacksRoute   = "/user/packs"
	userPackRoute    = "/user/packs/"
	authRoute        = "/auth"
	thumbnailRoute   = "/thumbnail"
)

var noAuthRoutes = []NoAuthRoute{
	{Path: baseRoute + authRoute, Method: http.MethodPost, PrefixMatch: false},
	{Path: baseRoute + publicPacksRoute, Method: http.MethodGet, PrefixMatch: false},
	{Path: baseRoute + thumbnailRoute, Method: http.MethodGet, PrefixMatch: false},
	{Path: "", Method: http.MethodOptions, PrefixMatch: true}, // preflight
}

func withCORS(h http.Handler) http.Handler {
	domain := config.Load().Domain()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domain)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
		if !strings.HasSuffix(r.URL.Path, thumbnailRoute) {
			w.Header().Set("Content-Type", "application/json")
		}
		h.ServeHTTP(w, r)
	})
}

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	api := http.NewServeMux()

	api.HandleFunc(authRoute, authHandler)

	api.HandleFunc(publicPacksRoute, publicPacksHandler)

	api.HandleFunc(userPackRoute, userPackHandler)
	api.HandleFunc(userPacksRoute, userPacksHandler)

	api.HandleFunc(thumbnailRoute, thumbnailHandler)

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
	id := strings.TrimPrefix(r.URL.Path, userPackRoute)
	if id == "" {
		http.Error(w, "Missing pack ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		deletePackHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userPacksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createPackHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := telegram.ParseAuth(r)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusBadRequest)
		return
	}

	jwt, err := SignID(req.ID)
	if err != nil {
		http.Error(w, "Failed to sign JWT", http.StatusBadRequest)
		return
	}

	// not sure how to store the domain yet
	// TODO FIX
	domain := config.Load().Domain()
	domain = strings.TrimPrefix(domain, "http://")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		// Uncomment in prod?
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
}
