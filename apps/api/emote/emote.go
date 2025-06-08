package emote

import (
	"fmt"
	"io"
	"net/http"
)

type Emote struct {
	URL        string   `json:"url"`
	Keywords   []string `json:"keywords"`
	Emoji_list []string `json:"emoji_list"`
}

func (emote Emote) Download() ([]byte, error) {
	resp, err := http.Get(emote.URL)
	if err != nil {
		return nil, fmt.Errorf("Failed to download %s: %w", emote.URL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad response for %s: %w", emote.URL, err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body for %s: %w", emote.URL, err)
	}

	return data, nil
}
