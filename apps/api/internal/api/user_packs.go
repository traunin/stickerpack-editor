package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/db"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/emote"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/resize"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
)

type DeletePackRequest struct {
	UserID   int64
	PackName string `json:"pack_name"`
}

type DeletePackResponse struct {
	Success bool `json:"success"`
}

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

func deletePackHandler(w http.ResponseWriter, r *http.Request) {
	req, mr := parseDeletePackRequest(w, r)
	if mr != nil {
		http.Error(w, mr.Error(), mr.status)
		return
	}

	pack, err := telegram.NewStickerPack(
		req.UserID,
		telegram.WithName(req.PackName),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pack.Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete sticker pack: %v", err), http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode(DeletePackResponse{Success: true})
}

func parseDeletePackRequest(
	w http.ResponseWriter,
	r *http.Request,
) (
	req *DeletePackRequest,
	mr *malformedRequest,
) {
	err := DecodeJSONBody(w, r, &req)
	if err != nil {
		if errors.As(err, &mr) {
			return
		}
		log.Printf("decoding error in parseDeletePacksRequest: %v", err)
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

func applyWatermark(title string, hasWatermark bool, cfg *config.Config) string {
	if hasWatermark {
		return fmt.Sprintf("%s by @%s", title, cfg.BotName())
	}
	return title
}

func emotesToStickers(emotes []emote.EmoteInput) []telegram.InputSticker {
	// TODO handle errors
	stickers := make([]telegram.InputSticker, len(emotes))
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
		stickers[i] = telegram.InputSticker{
			Sticker:   emoteData.File,
			Format:    format[emoteData.Animated],
			Keywords:  emote.Keywords(),
			EmojiList: emote.EmojiList(),
		}
	}

	return stickers
}

func packFromRequest(
	req *CreatePackRequest,
	cfg *config.Config,
) (
	*telegram.StickerPack,
	error,
) {
	watermarkTitle := applyWatermark(req.Title, req.HasWatermark, cfg)
	stickers := emotesToStickers(req.Emotes)
	return telegram.NewStickerPack(
		req.UserID,
		telegram.WithName(req.PackName),
		telegram.WithStickers(stickers),
		telegram.WithTitle(watermarkTitle),
		telegram.WithPublic(req.IsPublic),
	)
}

func createPackHandler(w http.ResponseWriter, r *http.Request) {
	req, mr := parseCreatePackRequest(w, r)
	if mr != nil {
		http.Error(w, mr.Error(), mr.status)
		return
	}

	cfg := config.Load()
	pack, err := packFromRequest(req, cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := pack.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	err = pack.UpdateThumbnailID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	storedPack := db.NewStoredPack(
		db.WithUserID(pack.UserID()),
		db.WithName(pack.Name()),
		db.WithTitle(pack.Title()),
		db.WithPublic(pack.IsPublic()),
		db.WithThumbnail(pack.ThumbnailID()),
	)
	err = cfg.DBConn().AddStickerpack(storedPack)
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
		log.Printf("decoding error in parseCreatePacksRequest: %v", err)
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
