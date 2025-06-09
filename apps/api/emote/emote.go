package emote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Emote struct {
	SevenTVID string   `json:"seventv_id"`
	Keywords  []string `json:"keywords"`
	EmojiList []string `json:"emoji_list"`
}

type EmoteData struct {
	Animated bool
	File     []byte
}

type SevenTVResponse struct {
	Animated bool `json:"animated"`
}

// I'm assuming there's always a gif and a png
var extensions = map[bool]string{
	true:  ".gif",
	false: ".png",
}

func (e Emote) Download() (EmoteData, error) {
	isAnimated, err := e.isAnimated()
	if err != nil {
		return EmoteData{}, fmt.Errorf("Failed to get data for %s: %w", e.SevenTVID, err)
	}

	extension := extensions[isAnimated]
	emoteURL := fmt.Sprintf("https://cdn.7tv.app/emote/%s/4x%s", e.SevenTVID, extension)

	resp, err := http.Get(emoteURL)
	if err != nil {
		return EmoteData{}, fmt.Errorf("Failed to send download request for %s: %w", emoteURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EmoteData{}, fmt.Errorf("Download response invalid for %s: %d", emoteURL, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmoteData{}, fmt.Errorf("Failed to read download for %s: %w", emoteURL, err)
	}

	emoteData := EmoteData{
		File:     data,
		Animated: isAnimated,
	}

	return emoteData, nil
}

func (e Emote) isAnimated() (bool, error) {
	// Currently using an old api, if it's deprecated...
	// We'll have to deal with GraphQL...

	requestURL := fmt.Sprintf("https://7tv.io/v3/emotes/%s", e.SevenTVID)
	resp, err := http.Get(requestURL)
	if err != nil {
		return false, fmt.Errorf("Failed to request %s: %w", requestURL, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return false, fmt.Errorf("Emote does not exist")
	case 400:
		return false, fmt.Errorf("Wrong ID")
	}

	var info SevenTVResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return false, fmt.Errorf("Failed to parse 7tv response for %s: %w", requestURL, err)
	}

	return info.Animated, nil
}
