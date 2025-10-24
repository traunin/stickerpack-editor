package emote

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/retrier"
)

var (
	retriesTenor    = config.Load().DownloadRetries()
	httpClientTenor = &http.Client{
		Timeout: 12 * time.Second,
	}
)

type tenorEmote struct {
	id        string
	keywords  []string
	emojiList []string
}

func (e *tenorEmote) Download(ctx context.Context) (EmoteData, error) {
	retryParams := &retrier.RetryParams{
		URL:     fmt.Sprintf("https://media.tenor.com/%s", e.id),
		Client:  httpClientTenor,
		Retries: retriesTenor,
	}
	data, err := retrier.Download(ctx, retryParams)
	if err != nil {
		return EmoteData{}, fmt.Errorf(
			"failed to download emote %s: %w", e.id, err,
		)
	}

	return EmoteData{
		File:     data,
		Animated: true,
	}, nil
}

func (e *tenorEmote) ID() string {
	return e.id
}

func (e *tenorEmote) Keywords() []string {
	return e.keywords
}

func (e *tenorEmote) EmojiList() []string {
	return e.emojiList
}

func (e *tenorEmote) String() string {
	return fmt.Sprintf("tenor:%s", e.id)
}
