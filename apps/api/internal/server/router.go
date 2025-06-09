package server

import (
	"net/http"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	api := http.NewServeMux()
	api.HandleFunc("/create-pack", createPackHandler)
	api.HandleFunc("/delete-pack", deletePackHandler)

	mux.Handle("/api/", http.StripPrefix("/api", api))
	return mux
}
