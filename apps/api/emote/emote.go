package emote

import (
	"fmt"
	"io"
	"net/http"
)

type Emote struct {
	SevenTVID   string   `json:"seventv_id"`
	Keywords  []string `json:"keywords"`
	EmojiList []string `json:"emoji_list"`
}

func (emote Emote) Download() ([]byte, error) {
	emoteURL := fmt.Sprintf("https://cdn.7tv.app/emote/%s/4x.png", emote.SevenTVID)
	resp, err := http.Get(emoteURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to download %s: %w", emoteURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad response for %s: %w", emoteURL, err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body for %s: %w", emoteURL, err)
	}

	return data, nil
}
