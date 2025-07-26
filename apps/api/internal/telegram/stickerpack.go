package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

type Sticker struct {
	Sticker   []byte
	Format    string
	EmojiList []string
	Keywords  []string
}

type StickerPack struct {
	UserID   int64
	Name     string
	Title    string
	Stickers []Sticker
}

type inputSticker struct {
	Sticker   string   `json:"sticker"`
	Format    string   `json:"format"`
	EmojiList []string `json:"emoji_list"`
	Keywords  []string `json:"keywords"`
}

func (pack StickerPack) Create() (string, error) {
	config := config.Load()
	botToken := config.TelegramToken
	botName := config.BotName
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/createNewStickerSet", botToken)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	IDstring := strconv.FormatInt(pack.UserID, 10)
	if err := writer.WriteField("user_id", IDstring); err != nil {
		return "", fmt.Errorf("failed to write user_id: %w", err)
	}
	validName := fmt.Sprintf("%s_by_%s", pack.Name, botName)
	if !isValidPackName(validName) {
		return "", errors.New("invalid stickerpack name")
	}
	if err := writer.WriteField("name", validName); err != nil {
		return "", fmt.Errorf("failed to write name: %w", err)
	}
	if err := writer.WriteField("title", pack.Title); err != nil {
		return "", fmt.Errorf("failed to write title: %w", err)
	}

	inputStickers := make([]inputSticker, len(pack.Stickers))
	for i, sticker := range pack.Stickers {
		inputStickers[i] = inputSticker{
			Sticker:   fmt.Sprintf("attach://sticker%d", i),
			EmojiList: sticker.EmojiList,
			Format:    sticker.Format,
			Keywords:  sticker.Keywords,
		}
	}

	jsonStickers, err := json.Marshal(inputStickers)
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %w", err)
	}
	writer.WriteField("stickers", string(jsonStickers))

	for i, sticker := range pack.Stickers {
		extension := ".png"
		if sticker.Format == "video" {
			extension = ".webm"
		}
		part, err := writer.CreateFormFile(fmt.Sprintf("sticker%d", i), fmt.Sprintf("sticker%d%s", i, extension))
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

func (pack StickerPack) Delete() error {
	config := config.Load()
	botToken := config.TelegramToken
	botName := config.BotName
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/deleteStickerSet", botToken)

	validName := fmt.Sprintf("%s_by_%s", pack.Name, botName)

	resp, err := http.PostForm(requestURL, url.Values{
		"name": {validName},
	})
	if err != nil {
		return fmt.Errorf("deleteStickerSet failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

func isValidPackName(name string) bool {
	// English letters and digits, underscores
	// <= 64 characters
	// no consequtive underscores
	re := regexp.MustCompile(`(?i)^[a-z][a-z0-9_]*$`)
	return len(name) <= 64 &&
		re.MatchString(name) &&
		!strings.Contains(name, "__")
}
