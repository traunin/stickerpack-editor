package telegram

type Sticker struct {
	FileID           string     `json:"file_id"`
	FileUniqueID     string     `json:"file_unique_id"`
	Type             string     `json:"type"`
	Width            int        `json:"width"`
	Height           int        `json:"height"`
	IsAnimated       bool       `json:"is_animated"`
	IsVideo          bool       `json:"is_video"`
	Thumbnail        *PhotoSize `json:"thumbnail,omitempty"`
	Emoji            string     `json:"emoji,omitempty"`
	SetName          string     `json:"set_name,omitempty"`
	PremiumAnimation *File      `json:"premium_animation,omitempty"`
	MaskPosition     any        `json:"mask_position,omitempty"`
	CustomEmojiID    string     `json:"custom_emoji_id,omitempty"`
	NeedsRepainting  bool       `json:"needs_repainting,omitempty"`
	FileSize         int        `json:"file_size,omitempty"`
}

type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
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
