package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/emote"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/resize"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
)

type CreatePackRequest struct {
	PackName     string             `json:"pack_name"`
	Title        string             `json:"title"`
	Emotes       []emote.EmoteInput `json:"emotes"`
	UserID       int64              `json:"user_id"`
	UseWatermark bool               `json:"use_watermark"`
}

type CreatePackResponse struct {
	PackURL string `json:"pack_url"`
}

var format = map[bool]string{
	true:  "video",
	false: "static",
}

func createPackHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stickers := make([]telegram.Sticker, len(req.Emotes))
	for i, input := range req.Emotes {
		emote, err := input.ToEmote()
		if err != nil {
			log.Printf("invalid emote input %s: %v", emote, err)
			continue
		}

		emoteData, err := emote.Download()
		if err != nil {
			log.Printf("failed downloading emote %s", emote)
			continue
		}
		err = resize.FitEmote(&emoteData)
		if err != nil {
			log.Printf("failed resizing emote %s: %v", emote, err)
			continue
		}
		stickers[i] = telegram.Sticker{
			Sticker:   emoteData.File,
			Format:    format[emoteData.Animated],
			Keywords:  emote.Keywords(),
			EmojiList: emote.EmojiList(),
		}
	}

	var title string
	if req.UseWatermark {
		title = fmt.Sprintf("%s by @%s", req.Title, config.Load().BotName())
	} else {
		title = req.Title
	}

	pack := telegram.StickerPack{
		UserID:   req.UserID,
		Name:     req.PackName,
		Title:    title,
		Stickers: stickers,
	}

	url, err := pack.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	err = config.Load().DBConn().AddStickerpack(req.UserID, title, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(CreatePackResponse{PackURL: url})
}

func parseRequest(r *http.Request) (*CreatePackRequest, error) {
	var req CreatePackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.New("invalid JSON schema")
	}

	if req.PackName == "" {
		return nil, errors.New("stickerpack name missing")
	}

	if req.Title == "" {
		return nil, errors.New("stickerpack title missing")
	}

	emoteCount := len(req.Emotes)
	if emoteCount == 0 {
		return nil, errors.New("no emotes in stickerpack")
	}

	return &req, nil
}
