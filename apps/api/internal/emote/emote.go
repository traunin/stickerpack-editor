package emote

import "fmt"

const maxKeywords = 20 - 5 // the tg limit is 20, we need 2, 5 is just to be safe
const maxEmojis = 20

type Emote interface {
	Download() (EmoteData, error)
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

func (e EmoteInput) ToEmote() (Emote, error) {
	if len(e.Keywords) > maxKeywords {
		return nil, fmt.Errorf("max %d keywords is supported", maxKeywords)
	}

	if len(e.EmojiList) > maxKeywords {
		return nil, fmt.Errorf("max %d emojis is supported", maxEmojis)
	}

	metaKeywords := append([]string{}, e.Source, e.ID)
	switch e.Source {
	case "7tv":
		return newSevenTVEmote(e.ID, metaKeywords, e.EmojiList)
	default:
		return nil, fmt.Errorf("unsupported source %s", e.Source)
	}
}
