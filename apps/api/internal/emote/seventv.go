package emote

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

const idLength = 26

type sevenTVEmote struct {
	id        string
	keywords  []string
	emojiList []string
}

type sevenTVResponse struct {
	Animated bool `json:"animated"`
}

func newSevenTVEmote(id string, keywords, emojiList []string) (sevenTVEmote, error) {
	// no checks for keywords and emojis since it's supposed to be checked in toEmote
	if !isValidId(id) {
		return sevenTVEmote{}, fmt.Errorf("id %s invalid", id)
	}

	return sevenTVEmote{
		id,
		keywords,
		emojiList,
	}, nil
}

func isValidId(id string) bool {
	return len(id) == idLength // best I can come up with right now
}

// I'm assuming there's always a gif and a png
var extensions = map[bool]string{
	true:  ".gif",
	false: ".png",
}

func (e sevenTVEmote) Download() (EmoteData, error) {
	isAnimated, err := e.isAnimated()
	if err != nil {
		return EmoteData{}, fmt.Errorf("failed to get data for %s: %w", e.id, err)
	}

	extension := extensions[isAnimated]
	emoteURL := fmt.Sprintf("https://cdn.7tv.app/emote/%s/4x%s", e.id, extension)

	data, err := downloadFile(emoteURL)
	if err != nil {
		return EmoteData{}, fmt.Errorf("failed to download emote %s: %w", e.id, err)
	}

	return EmoteData{
		File:     data,
		Animated: isAnimated,
	}, nil
}

func (e sevenTVEmote) ID() string {
	return e.id
}

func (e sevenTVEmote) Keywords() []string {
	return e.keywords
}

func (e sevenTVEmote) EmojiList() []string {
	return e.emojiList
}

func (e sevenTVEmote) String() string {
	return fmt.Sprintf("7tv:%s", e.id)
}

func (e sevenTVEmote) isAnimated() (bool, error) {
	retries := config.Load().DownloadRetries()

	for attempt := 0; attempt < retries; attempt++ {
		isAnimated, err := e.attemptIsAnimated()
		if err == nil {
			return isAnimated, nil
		}

		log.Printf("Checking animation status for %s failed %d/%d: %v", e.id, attempt+1, retries, err)

		// Don't sleep after the last attempt
		if attempt < retries-1 {
			sleepWithBackoff(attempt)
		}
	}

	return false, fmt.Errorf("failed to determine animation status for %s", e.id)
}

func (e sevenTVEmote) attemptIsAnimated() (bool, error) {
	// Currently using an old api, if it's deprecated...
	// We'll have to deal with GraphQL...
	requestURL := fmt.Sprintf("https://7tv.io/v3/emotes/%s", e.id)
	resp, err := http.Get(requestURL)
	if err != nil {
		return false, fmt.Errorf("failed to request %s: %w", requestURL, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return false, fmt.Errorf("emote does not exist")
	case http.StatusBadRequest:
		return false, fmt.Errorf("wrong ID")
	case http.StatusOK:
		var info sevenTVResponse
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return false, fmt.Errorf("failed to parse 7tv response for %s: %w", requestURL, err)
		}
		return info.Animated, nil
	default:
		return false, fmt.Errorf("request to %s returned unexpected status code %d", requestURL, resp.StatusCode)
	}
}
