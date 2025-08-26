package emote

import (
	"fmt"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

const (
	baseDelay = 200 * time.Millisecond
)

func downloadFile(url string) ([]byte, error) {
	retries := config.Load().DownloadRetries()

	for attempt := 0; attempt < retries; attempt++ {
		data, err := attemptDownload(url)
		if err == nil {
			return data, nil
		}

		log.Printf("Downloading %s failed %d/%d: %v", url, attempt+1, retries, err)

		// don't sleep after the last attempt
		if attempt < retries - 1 {
			sleepWithBackoff(attempt)
		}
	}

	return nil, fmt.Errorf("failed to download %s after %d attempts", url, retries)
}

func attemptDownload(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: %s returned %d", url, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func sleepWithBackoff(attempt int) {
	sleep := baseDelay * (1 << attempt)
	randomDelay := time.Duration(rand.Int63n(int64(sleep / 2)))
	time.Sleep(sleep + randomDelay)
}
