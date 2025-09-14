package emote

import (
	"fmt"
)

type tenorEmote struct {
	id        string
	keywords  []string
	emojiList []string
}

func newTenorEmote(id string, keywords, emojiList []string) (tenorEmote, error) {
	// no checks for keywords and emojis since it's supposed to be checked in toEmote
	return tenorEmote{
		id,
		keywords,
		emojiList,
	}, nil
}

func (e tenorEmote) Download() (EmoteData, error) {
	emoteURL := fmt.Sprintf("https://media.tenor.com/%s", e.id)

	data, err := downloadFile(emoteURL)
	if err != nil {
		return EmoteData{}, fmt.Errorf("failed to download emote %s: %w", e.id, err)
	}

	return EmoteData{
		File:     data,
		Animated: true,
	}, nil
}

func (e tenorEmote) ID() string {
	return e.id
}

func (e tenorEmote) Keywords() []string {
	return e.keywords
}

func (e tenorEmote) EmojiList() []string {
	return e.emojiList
}

func (e tenorEmote) String() string {
	return fmt.Sprintf("tenor:%s", e.id)
}
