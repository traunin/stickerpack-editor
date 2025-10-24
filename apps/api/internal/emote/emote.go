package emote

import (
	"context"
	"fmt"
)

const maxKeywords = 20 - 5 // the tg limit is 20, we need 2, 5 is just to be safe
const maxEmojis = 20

type Emote interface {
	Download(context.Context) (EmoteData, error)
	Keywords() []string
	EmojiList() []string
	ID() string
	String() string
}

type EmoteData struct {
	Animated bool
	File     []byte
}

type EmoteInput struct {
	Source    string   `json:"source"`
	ID        string   `json:"id"`
	Keywords  []string `json:"keywords"`
	EmojiList []string `json:"emoji_list"`
}

func (e *EmoteInput) ToEmote() (Emote, error) {
	if len(e.Keywords) > maxKeywords {
		return nil, fmt.Errorf("max %d keywords is supported", maxKeywords)
	}

	if len(e.EmojiList) > maxKeywords {
		return nil, fmt.Errorf("max %d emojis is supported", maxEmojis)
	}

	metaKeywords := append(append([]string{}, e.Keywords...), e.Source)
	switch e.Source {
	case "7tv":
        if !isValid7TVId(e.ID) {
            return nil, fmt.Errorf("id %s invalid", e.ID)
        }
        return &sevenTVEmote{e.ID, metaKeywords, e.EmojiList}, nil
	case "tenor":
		return &tenorEmote{e.ID, metaKeywords, e.EmojiList}, nil
	default:
		return nil, fmt.Errorf("unsupported source %s", e.Source)
	}
}
