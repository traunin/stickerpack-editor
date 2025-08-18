package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

type GetFileResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FilePath     string `json:"file_path"`
		FileSize     int    `json:"file_size"`
	} `json:"result"`
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()
	query := r.URL.Query()
	thumbnailID := query.Get("thumbnail_id")
	if thumbnailID == "" {
		http.Error(w, "missing thumbnail_id", http.StatusBadRequest)
		return
	}

	fileURL, err := downloadLink(cfg, thumbnailID)
	if err != nil {
		http.Error(w, "failed getting a download link", http.StatusBadGateway)
		return
	}

	var httpClient = &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(fileURL)
	if err != nil {
		http.Error(w, "failed to send download request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "download response invalid", http.StatusBadGateway)
		return
	}

	if strings.HasSuffix(fileURL, ".webm") {
		w.Header().Set("Content-Type", "video/webm")
	} else {
		w.Header().Set("Content-Type", "image/webp")
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "streaming failed", http.StatusInternalServerError)
	}
}

func downloadLink(cfg *config.Config, thumbnailID string) (string, error) {
	botToken := cfg.TelegramToken()
	reqURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getFile?file_id=%s",
		botToken,
		url.QueryEscape(thumbnailID),
	)
	resp, err := http.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("failed getting a download link: %w", err)
	}
	defer resp.Body.Close()

	fileResp, err := parseGetFileResponse(resp.Body)
	if err != nil {
		return "", err
	}

	if !fileResp.Ok {
		return "", fmt.Errorf("telegram API returned not ok")
	}
	if fileResp.Result.FilePath == "" {
		return "", fmt.Errorf("file_path is empty")
	}

	fileURL := fmt.Sprintf(
		"https://api.telegram.org/file/bot%s/%s",
		botToken,
		fileResp.Result.FilePath,
	)
	return fileURL, nil
}

func parseGetFileResponse(r io.Reader) (*GetFileResponse, error) {
	var fileResp GetFileResponse
	if err := json.NewDecoder(r).Decode(&fileResp); err != nil {
		return nil, fmt.Errorf("failed to decode getFile response: %w", err)
	}
	return &fileResp, nil
}
