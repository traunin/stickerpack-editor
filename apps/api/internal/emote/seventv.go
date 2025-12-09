package emote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/retrier"
)

const idLength = 26

var (
	retries7TV    = config.Load().DownloadRetries()
	httpClient7TV = &http.Client{
		Timeout: 10 * time.Second,
	}
)

type sevenTVEmote struct {
	id        string
	keywords  []string
	emojiList []string
}

type sevenTVResponse struct {
	Animated bool `json:"animated"`
}

func isValid7TVId(id string) bool {
	return len(id) == idLength // best I can come up with right now
}

// I'm assuming there's always a gif and a png
var extensions = map[bool]string{
	true:  ".gif",
	false: ".png",
}

func (e *sevenTVEmote) Download(ctx context.Context) (EmoteData, error) {
	isAnimated, err := e.isAnimated(ctx)
	if err != nil {
		return EmoteData{}, fmt.Errorf("failed to get data for %s: %w", e.id, err)
	}

	extension := extensions[isAnimated]
	url := fmt.Sprintf("https://cdn.7tv.app/emote/%s/4x%s", e.id, extension)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return EmoteData{}, fmt.Errorf("failed creating request: %w", err)
	}
	retryParams := &retrier.RetryParams{
		Request: req,
		Client:  httpClient7TV,
		Retries: retries7TV,
	}

	data, err := retrier.Download(retryParams)
	if err != nil {
		return EmoteData{}, fmt.Errorf("failed to download emote %s: %w", e.id, err)
	}

	return EmoteData{
		File:     data,
		Animated: isAnimated,
	}, nil
}

func (e *sevenTVEmote) ID() string {
	return e.id
}

func (e *sevenTVEmote) Keywords() []string {
	return e.keywords
}

func (e *sevenTVEmote) EmojiList() []string {
	return e.emojiList
}

func (e *sevenTVEmote) String() string {
	return fmt.Sprintf("7tv:%s", e.id)
}

func (e *sevenTVEmote) isAnimated(ctx context.Context) (bool, error) {
	// Currently using an old api, if it's deprecated...
	// We'll have to deal with GraphQL...
	url := fmt.Sprintf("https://7tv.io/v3/emotes/%s", e.id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed creating request: %w", err)
	}
	retryParams := &retrier.RetryParams{
		Request: req,
		Client:  httpClient7TV,
		Retries: retries7TV,
	}
	resp, err := retrier.RequestWithCallback(
		ctx, retryParams, animatedRespCallback,
	)
	if err != nil {
		return false, fmt.Errorf("failed to determine animated status: %v", err)
	}
	defer resp.Body.Close()

	var info sevenTVResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return false, fmt.Errorf("failed to parse 7tv response")
	}

	return info.Animated, nil
}

func animatedRespCallback(resp *http.Response) (bool, error) {
	switch resp.StatusCode {
	case http.StatusNotFound:
		return false, fmt.Errorf("emote does not exist")
	case http.StatusBadRequest:
		return false, fmt.Errorf("wrong ID")
	case http.StatusOK:
		return false, nil
	default:
		return true, fmt.Errorf(
			"request returned unexpected status code %d", resp.StatusCode,
		)
	}
}
