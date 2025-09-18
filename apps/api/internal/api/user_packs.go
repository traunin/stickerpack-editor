package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/db"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/emote"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/resize"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
	"golang.org/x/sync/errgroup"
)

type DeletePackResponse struct {
	Success bool `json:"success"`
}

type CreatePackRequest struct {
	UserID       int64              `json:"-"`
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

func deletePackHandler(w http.ResponseWriter, r *http.Request, name string) {
	userID, err := UserIDFromContext(r)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError)
		return
	}

	db := config.Load().DBConn()
	owned, err := db.IsPackOwner(name, userID)
	if !owned || err != nil {
		http.Error(w, "Can't confirm pack ownership", http.StatusUnauthorized)
		return
	}

	pack, err := telegram.NewStickerPack(
		userID,
		telegram.WithValidName(name),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pack.Delete()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete sticker pack: %v", err),
			http.StatusBadGateway,
		)
		return
	}

	err = db.DeletePack(name, userID)
	if err != nil {
		http.Error(w, "Failed to delete from db", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(DeletePackResponse{Success: true})
}

func applyWatermark(title string, hasWatermark bool, cfg *config.Config) string {
	if hasWatermark {
		return fmt.Sprintf("%s by @%s", title, cfg.BotName())
	}
	return title
}

func emotesToStickers(
	ctx context.Context,
	emotes []emote.EmoteInput,
	limit int,
	progress func(done, total int),
) ([]telegram.InputSticker, error) {
	stickers := make([]telegram.InputSticker, len(emotes))
	g, ctx := errgroup.WithContext(ctx)

	sem := make(chan struct{}, limit)
	var mu sync.Mutex
	completed := 0

	for i, input := range emotes {
		i, input := i, input

		g.Go(func() error {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return ctx.Err()
			}

			sticker, err := parseEmote(ctx, input)
			if err != nil {
				return err
			}
			stickers[i] = sticker

			mu.Lock()
			completed++
			progress(completed, len(emotes))
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return stickers, nil
}

func parseEmote(ctx context.Context, input emote.EmoteInput) (telegram.InputSticker, error) {
    emote, err := input.ToEmote()
    if err != nil {
        return telegram.InputSticker{}, err
    }
    
    select {
    case <-ctx.Done():
        return telegram.InputSticker{}, ctx.Err()
    default:
    }
    
    emoteData, err := emote.Download()
    if err != nil {
        return telegram.InputSticker{}, err
    }
    
    select {
    case <-ctx.Done():
        return telegram.InputSticker{}, ctx.Err()
    default:
    }
    
    if err := resize.FitEmote(&emoteData); err != nil {
        return telegram.InputSticker{}, err
    }

	return telegram.InputSticker{
		Sticker:   emoteData.File,
		Format:    format[emoteData.Animated],
		Keywords:  emote.Keywords(),
		EmojiList: emote.EmojiList(),
	}, nil
}

func packFromRequest(
	ctx context.Context,
	req *CreatePackRequest,
	cfg *config.Config,
	progress func(done, total int),
) (*telegram.StickerPack, error) {
	watermarkTitle := applyWatermark(req.Title, req.HasWatermark, cfg)
	stickers, err := emotesToStickers(ctx, req.Emotes, 2, progress)
	if err != nil {
		return nil, err
	}
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

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE unsupported", http.StatusInternalServerError)
		return
	}
	ctx := r.Context()

	progress := func(done, total int) {
		select {
		case <-ctx.Done():
			return
		default:
		}
		
		event := struct {
			Done  int `json:"done"`
			Total int `json:"total"`
		}{done, total}
		data, _ := json.Marshal(event)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	cfg := config.Load()
	pack, err := packFromRequest(ctx, req, cfg, progress)
	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: %q\n\n", err.Error())
		flusher.Flush()
		return
	}

	url, err := pack.Create()
	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: %q\n\n", err.Error())
		flusher.Flush()
		return
	}

	_ = pack.UpdateThumbnailID()
	storedPack := db.NewStoredPack(
		db.WithUserID(pack.UserID()),
		db.WithName(pack.Name()),
		db.WithTitle(pack.Title()),
		db.WithPublic(pack.IsPublic()),
		db.WithThumbnail(pack.ThumbnailID()),
	)
	createdPack, err := cfg.DBConn().AddStickerpack(storedPack)
	if err != nil {
		fmt.Fprintf(w, "event: error\ndata: %q\n\n", err.Error())
		flusher.Flush()
		return
	}

	result := struct {
		PackURL string           `json:"pack_url"`
		Pack    *db.PackResponse `json:"pack"`
	}{
		PackURL: url,
		Pack:    createdPack,
	}
	data, _ := json.Marshal(result)
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", data)
	flusher.Flush()
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

func getUserPacksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		http.Error(w, "page is not a number", http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil {
		http.Error(w, "page_size is not a number", http.StatusBadRequest)
		return
	}

	if page < 0 {
		http.Error(w, "page is less than zero", http.StatusBadRequest)
		return
	}
	if pageSize <= 0 {
		http.Error(w, "page_size has to be > 0", http.StatusBadRequest)
		return
	}

	userID, ctxErr := UserIDFromContext(r)
	if ctxErr != nil {
		http.Error(w, ctxErr.Error(), http.StatusBadRequest)
		return
	}

	db := config.Load().DBConn()
	packs, err := db.UserPacks(userID, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := db.UserPacksCount(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(GetPacksResponse{
		Packs: packs,
		Total: total,
	})
}

func nameExistsHandler(w http.ResponseWriter, r *http.Request, name string) {
	// race condition: two people create a pack with the same name
	// handled by telegram and UNIQUE(name), but it's not the solution
	// might store names of stickerpacks being parsed and fail early
	// but if the earlier pack with the same name fails...
	// this is just for the frontend to have a tick on the name input field
	cfg := config.Load()
	validName := telegram.ValidPackName(name)
	exists, err := cfg.DBConn().NameExists(validName)
	if err != nil {
		http.Error(w, "Database error: unable to check name", http.StatusInternalServerError)
		return
	}

	if exists {
		w.WriteHeader(http.StatusOK) // 200 name taken
	} else {
		http.NotFound(w, r) // 404 name available
	}
}
