package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/emote"
	"github.com/Traunin/stickerpack-editor/apps/api/resize"
	"github.com/Traunin/stickerpack-editor/apps/api/telegram"
)

type CreatePackRequest struct {
	PackName string        `json:"pack_name"`
	Title    string        `json:"title"`
	Emotes   []emote.Emote `json:"emotes"`
	UserID   string        `json:"userID"`
}

type CreatePackResponse struct {
	PackURL string `json:"pack_url"`
}

var format = map[bool]string{
	true:  "video",
	false: "static",
}

func createPackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not POST", http.StatusBadRequest)
		return
	}

	var req CreatePackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON schema", http.StatusBadRequest)
		return
	}

	if req.PackName == "" {
		http.Error(w, "Stickerpack name missing", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Stickerpack title missing", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "User ID missing", http.StatusBadRequest)
		return
	}

	emoteCount := len(req.Emotes)
	if emoteCount == 0 {
		http.Error(w, "No emotes in stickerpack", http.StatusBadRequest)
		return
	}

	stickers := make([]telegram.Sticker, emoteCount)

	for i, emote := range req.Emotes {
		emoteData, err := emote.Download()
		if err != nil {
			log.Printf("failed downloading emote %s", emote.SevenTVID)
			continue
		}
		err = resize.FitEmote(&emoteData)
		if err != nil {
			log.Printf("failed resizing emote %s: %v", emote.SevenTVID, err)
			continue
		}
		stickers[i] = telegram.Sticker{
			Sticker:   emoteData.File,
			Format:    format[emoteData.Animated],
			Keywords:  emote.Keywords,
			EmojiList: emote.EmojiList,
		}
	}

	pack := telegram.StickerPack{
		UserID:   req.UserID,
		Name:     req.PackName,
		Title:    req.Title,
		Stickers: stickers,
	}

	url, err := pack.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	json.NewEncoder(w).Encode(CreatePackResponse{PackURL: url})
}
