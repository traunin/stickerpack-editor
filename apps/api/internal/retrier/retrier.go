package retrier

import (
	"context"
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
	Request *http.Request
	Client  *http.Client
	Retries int
}

type RetryCallback func(*http.Response) error

func Download(params *RetryParams) ([]byte, error) {
	client := params.Client
	request := params.Request
	retries := params.Retries
	ctx := request.Context()
	url := request.URL

	for attempt := 1; attempt <= retries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("context cancelled: %w", err)
		}

		data, err := attemptDownload(request, client)
		if err == nil {
			return data, nil
		}

		log.Printf(
			"downloading %s failed (%d/%d): %v", url, attempt, retries, err,
		)

		// don't sleep after the last attempt
		if attempt < retries {
			if err := sleepWithBackoff(ctx, attempt-1); err != nil {
				return nil, err
			}
		}
	}

	err := fmt.Errorf("failed to download %s after %d attempts", url, retries)
	return nil, err
}

func attemptDownload(
	request *http.Request,
	client *http.Client,
) ([]byte, error) {
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if statusCode := resp.StatusCode; statusCode != http.StatusOK {
		err := fmt.Errorf("download failed: %s returned %d", request.URL, statusCode)
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

// RequestWithCallback performs an HTTP GET request with retries and exponential backoff.
// The callback is invoked with each response. If the callback returns true, the response
// is returned to the caller, and it is the caller’s responsibility to close resp.Body.
// Otherwise, the response body is automatically closed before the next attempt.
func RequestWithCallback(
	ctx context.Context,
	params *RetryParams,
	callback RetryCallback,
) (*http.Response, error) {
	client := params.Client
	request := params.Request
	retries := params.Retries
	url := request.URL

	for attempt := 1; attempt <= retries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("context cancelled: %w", err)
		}

		resp, err := attemptRequest(request, client)
		if err != nil {
			log.Printf(
				"request of %s failed (%d/%d): %v", url, attempt, retries, err,
			)
			if attempt < retries {
				if err := sleepWithBackoff(ctx, attempt-1); err != nil {
					return nil, err
				}
			}
			continue
		}

		if err := invokeCallback(resp, callback); err == nil {
			return resp, nil
		}

		// don't sleep after the last attempt
		if attempt < retries {
			if err := sleepWithBackoff(ctx, attempt-1); err != nil {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("no valid response in %d requests", retries)
}

func attemptRequest(
	req *http.Request,
	client *http.Client,
) (*http.Response, error) {
	return client.Do(req)
}

func invokeCallback(resp *http.Response, callback RetryCallback) error {
	var callbackErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				resp.Body.Close()
				callbackErr = fmt.Errorf("callback panicked: %v", r)
			}
		}()
		callbackErr = callback(resp)
	}()

	if callbackErr != nil {
		resp.Body.Close()
	}

	return callbackErr
}

func sleepWithBackoff(ctx context.Context, attempt int) error {
	sleep := baseDelay * (1 << attempt)
	randomDelay := time.Duration(rand.Int63n(int64(sleep / 2)))
	totalSleep := sleep + randomDelay

	timer := time.NewTimer(totalSleep)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
