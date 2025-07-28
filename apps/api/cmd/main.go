package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/api"
)

func init() {
	_, err := exec.LookPath("ffmpeg")

	if err != nil {
		log.Fatal("no ffmpeg found on PATH, consider reloading or install if you haven't.")
	}
}

func main() {
	config := config.Load()
	router := api.SetupRouter()

	addr := ":" + config.Port()
	log.Printf("Server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}

	log.Println("Exiting...")
}
