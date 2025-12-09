package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/retrier"
)

const (
	cacheControl = "public, max-age=36000"
)

var (
	mediaRetries   = config.Load().DownloadRetries()
	httpClient     = &http.Client{Timeout: 15 * time.Second}
	fileInfoCache  = cache.New(55*time.Minute, 20*time.Minute)
	detectionLocks sync.Map
)

type CachedFileInfo struct {
	URL         string
	ContentType string
}

type StreamContext struct {
	Writer  http.ResponseWriter
	Data    []byte
	FileID  string
	FileURL string
}

type FileStreamRequest struct {
	Writer   http.ResponseWriter
	FileInfo *CachedFileInfo
	FileID   string
}

type DetectedStreamContext struct {
	Writer      http.ResponseWriter
	Data        []byte
	ContentType string
}

type GetFileResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FilePath     string `json:"file_path"`
		FileSize     int    `json:"file_size"`
	} `json:"result"`
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()
	ctx := r.Context()

	fileID, err := extractFileID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileInfo, err := getCachedOrFetchFileInfo(ctx, cfg, fileID)
	if err != nil {
		log.Printf("Error fetching file info for %s: %v\n", fileID, err)
		http.Error(w, "failed getting a download link", http.StatusBadGateway)
		return
	}

	req := FileStreamRequest{
		Writer:   w,
		FileInfo: fileInfo,
		FileID:   fileID,
	}
	if err := streamFileAndMaybeDetect(ctx, req); err != nil {
		if !isClientDisconnect(err) {
			log.Printf("Error streaming file %s: %v\n", fileID, err)
		}
		return
	}
}

func getCachedOrFetchFileInfo(
	ctx context.Context,
	cfg *config.Config,
	fileID string,
) (*CachedFileInfo, error) {
	if info, found := fileInfoCache.Get(fileID); found {
		return info.(*CachedFileInfo), nil
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	fileURL, err := downloadLink(ctx, cfg, fileID)
	if err != nil {
		return nil, err
	}

	info := &CachedFileInfo{
		URL:         fileURL,
		ContentType: "",
	}
	fileInfoCache.Set(fileID, info, cache.DefaultExpiration)

	return info, nil
}

func extractFileID(r *http.Request) (string, error) {
	query := r.URL.Query()
	fileID := query.Get("file_id")
	if fileID == "" {
		return "", fmt.Errorf("missing file_id")
	}
	return fileID, nil
}

func streamFileAndMaybeDetect(ctx context.Context, req FileStreamRequest) error {
	data, err := downloadFile(ctx, req.FileInfo.URL)
	if err != nil {
		return err
	}

	// Content-Type already cached
	if contentType := req.FileInfo.ContentType; contentType != "" {
		return streamWithCachedType(req.Writer, data, contentType)
	}

	// detect content type
	streamCtx := StreamContext{
		Writer:  req.Writer,
		Data:    data,
		FileID:  req.FileID,
		FileURL: req.FileInfo.URL,
	}
	return detectAndStream(streamCtx)
}

func streamWithCachedType(
	w http.ResponseWriter,
	data []byte,
	contentType string,
) error {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", cacheControl)

	if _, err := w.Write(data); err != nil {
		if isClientDisconnect(err) {
			return nil
		}
		return fmt.Errorf("streaming failed: %w", err)
	}

	return nil
}

func detectAndStream(ctx StreamContext) error {
	mu := acquireDetectionLock(ctx.FileID)
	mu.Lock()
	defer mu.Unlock()
	defer releaseDetectionLock(ctx.FileID)

	if contentType := getCachedContentType(ctx.FileID); contentType != "" {
		return streamWithCachedType(ctx.Writer, ctx.Data, contentType)
	}

	contentType := detectContentTypeFromStream(ctx.Data)
	updateCachedContentType(ctx.FileID, ctx.FileURL, contentType)

	streamCtx := DetectedStreamContext{
		Writer:      ctx.Writer,
		Data:        ctx.Data,
		ContentType: contentType,
	}
	return streamWithDetectedType(streamCtx)
}

func acquireDetectionLock(fileID string) *sync.Mutex {
	lockKey := "detect:" + fileID
	actualLock, _ := detectionLocks.LoadOrStore(lockKey, &sync.Mutex{})
	return actualLock.(*sync.Mutex)
}

func releaseDetectionLock(fileID string) {
	detectionLocks.Delete("detect:" + fileID)
}

func getCachedContentType(fileID string) string {
	if cached, found := fileInfoCache.Get(fileID); found {
		info := cached.(*CachedFileInfo)
		return info.ContentType
	}
	return ""
}

func detectContentTypeFromStream(data []byte) string {
	contentType := http.DetectContentType(data)
	return contentType
}

func updateCachedContentType(fileID, fileURL, contentType string) {
	updatedInfo := &CachedFileInfo{
		URL:         fileURL,
		ContentType: contentType,
	}
	fileInfoCache.Set(fileID, updatedInfo, cache.DefaultExpiration)
}

func streamWithDetectedType(ctx DetectedStreamContext) error {
	ctx.Writer.Header().Set("Content-Type", ctx.ContentType)
	ctx.Writer.Header().Set("Cache-Control", cacheControl)

	if _, err := ctx.Writer.Write(ctx.Data); err != nil {
		if isClientDisconnect(err) {
			return nil
		}
		return fmt.Errorf("streaming failed: %w", err)
	}

	return nil
}

func isClientDisconnect(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "connection timed out") ||
		strings.Contains(errStr, "client disconnected")
}

func downloadFile(ctx context.Context, fileURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}
	params := &retrier.RetryParams{
		Request: req,
		Client:  httpClient,
		Retries: mediaRetries,
	}
	return retrier.Download(params)
}

func downloadLink(
	ctx context.Context,
	cfg *config.Config,
	fileID string,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	botToken := cfg.TelegramToken()
	reqURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getFile?file_id=%s",
		botToken,
		url.QueryEscape(fileID),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed creating request: %w", err)
	}
	params := &retrier.RetryParams{
		Request: req,
		Client:  httpClient,
		Retries: mediaRetries,
	}
	resp, err := retrier.RequestWithCallback(ctx, params, downloadLinkCallback)
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

func downloadLinkCallback(resp *http.Response) (bool, error) {
	switch resp.StatusCode {
	case http.StatusOK:
		return false, nil
	default:
		return true,fmt.Errorf(
			"tg request for download link status code %d", resp.StatusCode,
		)
	}
}

func parseGetFileResponse(r io.Reader) (*GetFileResponse, error) {
	var fileResp GetFileResponse
	if err := json.NewDecoder(r).Decode(&fileResp); err != nil {
		return nil, fmt.Errorf("failed to decode getFile response: %w", err)
	}
	return &fileResp, nil
}
