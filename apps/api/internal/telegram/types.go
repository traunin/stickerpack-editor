package telegram

type Sticker struct {
	FileID           string      `json:"file_id"`
	FileUniqueID     string      `json:"file_unique_id"`
	Type             string      `json:"type"`
	Width            int         `json:"width"`
	Height           int         `json:"height"`
	IsAnimated       bool        `json:"is_animated"`
	IsVideo          bool        `json:"is_video"`
	Thumbnail        *PhotoSize  `json:"thumbnail,omitempty"`
	Emoji            string      `json:"emoji,omitempty"`
	SetName          string      `json:"set_name,omitempty"`
	PremiumAnimation *File       `json:"premium_animation,omitempty"`
	MaskPosition     interface{} `json:"mask_position,omitempty"`
	CustomEmojiID    string      `json:"custom_emoji_id,omitempty"`
	NeedsRepainting  bool        `json:"needs_repainting,omitempty"`
	FileSize         int         `json:"file_size,omitempty"`
}

type File struct {
	File_id        string `json:"file_id"`
	File_unique_id string `json:"file_unique_id"`
	File_size      int    `json:"file_size,omitempty"`
	File_path      string `json:"file_path,omitempty"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

type StickerSet struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	StickerType string     `json:"sticker_type"`
	Stickers    []Sticker  `json:"stickers"`
	Thumbnail   *PhotoSize `json:"thumbnail,omitempty"`
}
