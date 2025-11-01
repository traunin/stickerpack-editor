package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/api"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

func init() {
	_, err := exec.LookPath("ffmpeg")

	if err != nil {
		log.Fatalf("no ffmpeg on PATH: %v", err)
	}
}

func main() {
	cfg := config.Load()
	handler := api.SetupHandler()
	addr := ":" + cfg.Port()
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}

		log.Println("stopped server")
	}()
	log.Printf("listening on port %s", addr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("shutdown complete")
}
