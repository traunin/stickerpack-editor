package api

import (
	"net/http"
	"strings"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
)

func createSessionHandler(w http.ResponseWriter, r *http.Request) {
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
	domain := config.Load().Domain()
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   true, // works with https
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
}

func deleteSessionHandler(w http.ResponseWriter, _ *http.Request) {
	domain := config.Load().Domain()
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		MaxAge:   -1,
	})
}
