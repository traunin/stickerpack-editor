package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/api"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

func init() {
	_, err := exec.LookPath("ffmpeg")

	if err != nil {
		log.Fatal("no ffmpeg on PATH: %w", err)
	}
}

func main() {
	cfg := config.Load()
	handler := api.SetupHandler()

	addr := ":" + cfg.Port()
	log.Printf("listening on port %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}

	log.Println("stopping...")
}
