package retrier

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	baseDelay = 200 * time.Millisecond
)

type RetryParams struct {
	URL     string
	Client  *http.Client
	Retries int
}

type RetryCallback func(*http.Response) error

func Download(params *RetryParams) ([]byte, error) {
	client := params.Client
	url := params.URL
	retries := params.Retries

	for attempt := 1; attempt <= retries; attempt++ {
		data, err := attemptDownload(url, client)
		if err == nil {
			return data, nil
		}

		log.Printf(
			"downloading %s failed (%d/%d): %v", url, attempt, retries, err,
		)

		// don't sleep after the last attempt
		if attempt < retries {
			sleepWithBackoff(attempt - 1)
		}
	}

	err := fmt.Errorf("failed to download %s after %d attempts", url, retries)
	return nil, err
}

func attemptDownload(url string, client *http.Client) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if statusCode := resp.StatusCode; statusCode != http.StatusOK {
		err := fmt.Errorf("download failed: %s returned %d", url, statusCode)
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

// RequestWithCallback performs an HTTP GET request with retries and exponential backoff.
// The callback is invoked with each response. If the callback returns true, the response
// is returned to the caller, and it is the caller’s responsibility to close resp.Body.
// Otherwise, the response body is automatically closed before the next attempt.
func RequestWithCallback(
	params *RetryParams,
	callback RetryCallback,
) (*http.Response, error) {
	client := params.Client
	url := params.URL
	retries := params.Retries

	for attempt := 1; attempt <= retries; attempt++ {
		resp, err := client.Get(url)
		if err != nil {
			log.Printf(
				"request of %s failed (%d/%d): %v", url, attempt, retries, err,
			)
			continue
		}

		// if callback panics
		var callbackErr error = nil
		func() {
			defer func() {
				if r := recover(); r != nil {
					resp.Body.Close()
					callbackErr = fmt.Errorf("callback panicked: %v", r)
				}
			}()
			callbackErr = callback(resp)
		}()

		if callbackErr == nil {
			return resp, nil
		}
		resp.Body.Close()

		// don't sleep after the last attempt
		if attempt < retries {
			sleepWithBackoff(attempt - 1)
		}
	}

	return nil, fmt.Errorf("no valid response in %d requests", retries)
}

func Request(params *RetryParams) (*http.Response, error) {
	return RequestWithCallback(params, func(r *http.Response) error {
		return nil
	})
}

func sleepWithBackoff(attempt int) {
	sleep := baseDelay * (1 << attempt)
	randomDelay := time.Duration(rand.Int63n(int64(sleep / 2)))
	time.Sleep(sleep + randomDelay)
}
