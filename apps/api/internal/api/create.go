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
	UserID       int64
	PackName     string             `json:"pack_name"`
	Title        string             `json:"title"`
	Emotes       []emote.EmoteInput `json:"emotes"`
	IsPublic     bool               `json:"is_public"`
	HasWatermark bool               `json:"has_watermark"`
}

type CreatePackResponse struct {
	PackURL string `json:"pack_url"`
}

var format = map[bool]string{
	true:  "video",
	false: "static",
}

func applyWatermark(title string, hasWatermark bool, cfg *config.Config) string {
	if hasWatermark {
		return fmt.Sprintf("%s by @%s", title, cfg.BotName())
	}
	return title
}

func emotesToStickers(emotes []emote.EmoteInput) []telegram.Sticker {
	// TODO handle errors
	stickers := make([]telegram.Sticker, 0, len(emotes))
	for i, input := range emotes {
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
	
	return stickers
}

func createPackHandler(w http.ResponseWriter, r *http.Request) {
	req, mr := parseCreatePackRequest(w, r)
	if mr != nil {
		http.Error(w, mr.Error(), mr.status)
		return
	}

	stickers := emotesToStickers(req.Emotes)

	cfg := config.Load()
	title := applyWatermark(req.Title, req.HasWatermark, cfg)

	pack := telegram.StickerPack{
		UserID:   req.UserID,
		Name:     req.PackName,
		Title:    title,
		Stickers: stickers,
	}

	url, err := pack.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	err = cfg.DBConn().AddStickerpack(req.UserID, title, req.IsPublic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CreatePackResponse{PackURL: url})
}

func parseCreatePackRequest(
	w http.ResponseWriter,
	r *http.Request,
) (
	req *CreatePackRequest,
	mr *malformedRequest,
) {
	err := DecodeJSONBody(w, r, &req)
	if err != nil {
		if errors.As(err, &mr) {
			return
		}
		log.Printf("decoding error in parseGetPacksRequest: %v", err)
		mr = &malformedRequest{
			status: http.StatusInternalServerError,
			msg:    "unable to decode request",
		}
		return
	}

	if req.PackName == "" {
		mr = &malformedRequest{
			status: http.StatusBadRequest,
			msg:    "pack name is empty",
		}
		return
	}

	if req.Title == "" {
		mr = &malformedRequest{
			status: http.StatusBadRequest,
			msg:    "pack title is empty",
		}
		return
	}

	emoteCount := len(req.Emotes)
	if emoteCount == 0 {
		mr = &malformedRequest{
			status: http.StatusBadRequest,
			msg:    "no emotes in pack",
		}
		return
	}

	userID, ctxErr := UserIDFromContext(r)
	if ctxErr != nil {
		mr = &malformedRequest{
			status: http.StatusBadRequest,
			msg:    ctxErr.Error(),
		}
		return
	}
	req.UserID = userID

	return
}
