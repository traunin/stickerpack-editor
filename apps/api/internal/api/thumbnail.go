package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

const contentTypeBufferSize = 512
var httpClient = &http.Client{Timeout: 15 * time.Second}

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
	
	thumbnailID, err := extractThumbnailID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileURL, err := downloadLink(cfg, thumbnailID)
	if err != nil {
		http.Error(w, "failed getting a download link", http.StatusBadGateway)
		return
	}

	if err := streamFile(w, fileURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
}

func extractThumbnailID(r *http.Request) (string, error) {
	query := r.URL.Query()
	thumbnailID := query.Get("thumbnail_id")
	if thumbnailID == "" {
		return "", fmt.Errorf("missing thumbnail_id")
	}
	return thumbnailID, nil
}

func streamFile(w http.ResponseWriter, fileURL string) error {
	resp, err := downloadFile(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read some bits to detect file type
	buffer := make([]byte, contentTypeBufferSize)
	n, err := resp.Body.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read response")
	}

	contentType := detectContentType(buffer, n)
	w.Header().Set("Content-Type", contentType)

	// write buffer
	if _, err := w.Write(buffer[:n]); err != nil {
		return fmt.Errorf("failed to write buffer")
	}

	// Copy the rest
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("streaming failed")
	}

	return nil
}

func detectContentType(buffer []byte, n int) string {
	return http.DetectContentType(buffer[:n])
}

func downloadFile(fileURL string) (*http.Response, error) {
	resp, err := httpClient.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send download request")
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download response invalid")
	}

	return resp, nil
}

func downloadLink(cfg *config.Config, thumbnailID string) (string, error) {
	botToken := cfg.TelegramToken()
	reqURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getFile?file_id=%s",
		botToken,
		url.QueryEscape(thumbnailID),
	)
	resp, err := httpClient.Get(reqURL)
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
