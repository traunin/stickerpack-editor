package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/config"
)

type Sticker struct {
	Sticker    []byte
	Format     string
	Emoji_list []string
	Keywords   []string
}

type StickerPack struct {
	UserID   string
	Name     string
	Title    string
	Stickers []Sticker
}

type inputSticker struct {
	Sticker    string   `json:"sticker"`
	Format     string   `json:"format"`
	Emoji_list []string `json:"emoji_list"`
	Keywords   []string `json:"keywords"`
}

func (pack StickerPack) Create() (string, error) {
	config := config.Load()
	botToken := config.TelegramToken
	botName := config.BotName
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/createNewStickerSet", botToken)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("user_id", pack.UserID)
	validName := fmt.Sprintf("%s_by_%s", pack.Name, botName)
	writer.WriteField("name", validName)
	writer.WriteField("title", pack.Title)

	inputStickers := make([]inputSticker, len(pack.Stickers))
	for i, sticker := range pack.Stickers {
		inputStickers[i] = inputSticker{
			Sticker:    fmt.Sprintf("attach://sticker%d", i),
			Emoji_list: sticker.Emoji_list,
			Format:     sticker.Format,
			Keywords:   sticker.Keywords,
		}
	}

	jsonStickers, err := json.Marshal(inputStickers)
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %w", err)
	}
	writer.WriteField("stickers", string(jsonStickers))

	for i, sticker := range pack.Stickers {
		part, err := writer.CreateFormFile(fmt.Sprintf("sticker%d", i), fmt.Sprintf("sticker%d.webp", i))
		if err != nil {
			return "", fmt.Errorf("failed writing to request: %w", err)
		}
		if _, err := part.Write(sticker.Sticker); err != nil {
			return "", fmt.Errorf("failed attaching image: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	resp, err := http.Post(requestURL, writer.FormDataContentType(), &buf)
	if err != nil {
		return "", fmt.Errorf("createNewStickerSet failed: %w", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("telegram API error: %s", string(body))
	}

	stickerPackURL := fmt.Sprintf("https://t.me/addstickers/%s", validName)
	return stickerPackURL, nil
}
