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

type StickerSetResponse struct {
	telegram.StickerSet
	IsPublic bool `json:"is_public"`
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

type EditPackRequest struct {
	UserID          int64                   `json:"-"`
	PackName        string                  `json:"-"`
	UpdatedTitle    *string                 `json:"updated_title,omitempty"`
	UpdatedIsPublic *bool                   `json:"updated_is_public,omitempty"`
	DeletedStickers []string                `json:"deleted_stickers"`
	AddedStickers   []emote.EmoteInput      `json:"added_stickers"`
	EmojiUpdates    []StickerEmojiUpdate    `json:"emoji_updates"`
	PositionUpdates []StickerPositionUpdate `json:"position_updates"`
}

type StickerEmojiUpdate struct {
	ID     string   `json:"id"`
	Emojis []string `json:"emojis"`
}

type StickerPositionUpdate struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
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

func getPackHandler(w http.ResponseWriter, r *http.Request, name string) {
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

	set, err := pack.Fetch(r.Context())
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to fetch stickerpack: %v", err),
			http.StatusBadGateway,
		)
		return
	}

	isPublic, err := db.IsPackPublic(name)
	if err != nil {
		http.Error(w, "Failed to query publicity", http.StatusInternalServerError)
		return
	}

	resp := StickerSetResponse{
		StickerSet: *set,
		IsPublic:   isPublic,
	}

	json.NewEncoder(w).Encode(resp)
}

func applyWatermark(title string, hasWatermark bool, cfg *config.Config) string {
	if hasWatermark {
		return fmt.Sprintf("%s by @%s", title, cfg.BotName())
	}
	return title
}

type CreatePackJobHandler struct {
	cfg *config.Config
	req *CreatePackRequest
}

func NewCreatePackJobHandler(
	cfg *config.Config,
	req *CreatePackRequest,
) *CreatePackJobHandler {
	return &CreatePackJobHandler{
		cfg: cfg,
		req: req,
	}
}

func (h *CreatePackJobHandler) GetJobType() string {
	return "create_stickerpack"
}

func createPackHandler(w http.ResponseWriter, r *http.Request) {
	req, mr := parseCreatePackRequest(w, r)
	if mr != nil {
		http.Error(w, mr.Error(), mr.status)
		return
	}

	cfg := config.Load()
	handler := NewCreatePackJobHandler(cfg, req)

	jobID, err := jobQueue.Enqueue(handler, r)
	if err != nil {
		msg := fmt.Sprintf("Failed to enqueue job: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"job_id": jobID,
		"status": "queued",
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
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
			mu.Unlock()
			progress(completed, len(emotes))
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return stickers, nil
}

func parseEmote(
	ctx context.Context,
	input emote.EmoteInput,
) (telegram.InputSticker, error) {
	emote, err := input.ToEmote()
	if err != nil {
		return telegram.InputSticker{}, err
	}

	if err := ctx.Err(); err != nil {
		return telegram.InputSticker{}, err
	}

	emoteData, err := emote.Download(ctx)
	if err != nil {
		return telegram.InputSticker{}, err
	}

	if err := ctx.Err(); err != nil {
		return telegram.InputSticker{}, err
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

func (h *CreatePackJobHandler) Handle(
	ctx context.Context,
	r *http.Request,
	progress func(done, total int, message string),
) (any, error) {
	req := h.req
	steps := 3 + len(req.Emotes)
	currentStep := 0

	progress(currentStep, steps, "Processing emotes")
	stickers, err := emotesToStickers(
		ctx,
		req.Emotes,
		2,
		func(done, total int) {
			currentStep = steps - total + done
			progress(
				currentStep,
				steps,
				fmt.Sprintf("Processing emotes (%d/%d)", done, total),
			)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to process emotes: %w", err)
	}

	progress(currentStep, steps, "Creating stickerpack")
	currentStep++
	watermarkTitle := applyWatermark(req.Title, req.HasWatermark, h.cfg)
	pack, err := telegram.NewStickerPack(
		req.UserID,
		telegram.WithName(req.PackName),
		telegram.WithStickers(stickers),
		telegram.WithTitle(watermarkTitle),
		telegram.WithPublic(req.IsPublic),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bundle stickerpack: %w", err)
	}

	url, err := pack.Create()
	if err != nil {
		return nil, fmt.Errorf("telegram error: %w", err)
	}

	progress(currentStep, steps, "Saving to database")
	_ = pack.UpdateThumbnailID()
	storedPack := db.NewStoredPack(
		db.WithUserID(pack.UserID()),
		db.WithName(pack.Name()),
		db.WithTitle(pack.Title()),
		db.WithPublic(pack.IsPublic()),
		db.WithThumbnail(pack.ThumbnailID()),
	)

	createdPack, err := h.cfg.DBConn().AddStickerpack(storedPack)
	if err != nil {
		return nil, fmt.Errorf("failed to save pack to database: %w", err)
	}

	return struct {
		PackURL string           `json:"pack_url"`
		Pack    *db.PackResponse `json:"pack"`
	}{
		PackURL: url,
		Pack:    createdPack,
	}, nil
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

	resp, err := userPacksPreviews(r.Context(), userID, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func userPacksPreviews(
	ctx context.Context,
	userID int64,
	page,
	pageSize int,
) (*GetPacksResponse, error) {
	db := config.Load().DBConn()
	packs, err := db.UserPacks(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	previews := make([]telegram.PackPreview, len(packs))
	for i := range packs {
		preview, err := telegram.FetchPackPreview(ctx, packs[i].Name)
		if err != nil {
			return nil, err
		}
		previews[i] = *preview
	}

	total, err := db.UserPacksCount(userID)
	if err != nil {
		return nil, err
	}

	return &GetPacksResponse{
		Packs: previews,
		Total: total,
	}, nil
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
		http.Error(
			w,
			"Database error: unable to check name",
			http.StatusInternalServerError,
		)
		return
	}

	if exists {
		w.WriteHeader(http.StatusOK) // 200 name taken
	} else {
		http.NotFound(w, r) // 404 name available
	}
}

type EditPackJobHandler struct {
	cfg *config.Config
	req *EditPackRequest
}

func NewEditPackJobHandler(
	cfg *config.Config,
	req *EditPackRequest,
) *EditPackJobHandler {
	return &EditPackJobHandler{
		cfg: cfg,
		req: req,
	}
}

func (h *EditPackJobHandler) GetJobType() string {
	return "edit_stickerpack"
}

func editPackHandler(w http.ResponseWriter, r *http.Request, name string) {
	req, mr := parseEditPackRequest(w, r, name)
	if mr != nil {
		http.Error(w, mr.Error(), mr.status)
		return
	}

	cfg := config.Load()
	handler := NewEditPackJobHandler(cfg, req)

	jobID, err := jobQueue.Enqueue(handler, r)
	if err != nil {
		msg := fmt.Sprintf("Failed to enqueue job: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"job_id": jobID,
		"status": "queued",
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func parseEditPackRequest(
	w http.ResponseWriter,
	r *http.Request,
	name string,
) (
	req *EditPackRequest,
	mr *malformedRequest,
) {
	err := DecodeJSONBody(w, r, &req)
	if err != nil {
		if errors.As(err, &mr) {
			return
		}
		log.Printf("decoding error in parseEditPacksRequest: %v", err)
		mr = &malformedRequest{
			status: http.StatusInternalServerError,
			msg:    "unable to decode request",
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
	req.PackName = name
	return
}

type editProgress struct {
	current int
	total   int
	notify  func(done, total int, message string)
}

func (p *editProgress) update(message string) {
	p.current++
	p.notify(p.current, p.total, message)
}

func (p *editProgress) setMessage(message string) {
	p.notify(p.current, p.total, message)
}

type editResponse struct {
	Pack telegram.PackPreview `json:"pack"`
}

func (h *EditPackJobHandler) Handle(
	ctx context.Context,
	r *http.Request,
	progress func(done, total int, message string),
) (any, error) {
	req := h.req
	name := req.PackName
	prog := &editProgress{
		current: 0,
		total:   calculateEditSteps(req),
		notify:  progress,
	}
	pack, err := telegram.NewStickerPack(
		req.UserID,
		telegram.WithValidName(name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bundle stickerpack: %w", err)
	}

	prog.setMessage("Starting pack edit")
	if err := editDeleteStage(req.DeletedStickers, prog); err != nil {
		return nil, fmt.Errorf("failed to delete stickers: %w", err)
	}
	if err := editEmojiStage(req.EmojiUpdates, prog); err != nil {
		return nil, fmt.Errorf("failed to update emojis: %w", err)
	}
	if err := editUpdateIsPublicStage(req, name, prog); err != nil {
		return nil, fmt.Errorf("failed to update IsPublic: %w", err)
	}
	if err := editUpdateTitleStage(req, pack, prog); err != nil {
		return nil, fmt.Errorf("failed to update title: %w", err)
	}
	if err := editAddStage(ctx, pack, req.AddedStickers, prog); err != nil {
		return nil, fmt.Errorf("failed to add stickers: %w", err)
	}
	if err := editPositionStage(req.PositionUpdates, prog); err != nil {
		return nil, fmt.Errorf("failed to update positions: %w", err)
	}

	preview, err := telegram.FetchPackPreview(ctx, req.PackName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch edited pack: %w", err)
	}

	return editResponse{Pack: *preview}, nil
}

func calculateEditSteps(req *EditPackRequest) int {
	totalSteps := 0
	totalSteps += len(req.DeletedStickers)
	totalSteps += len(req.AddedStickers) * 2 // process and add
	totalSteps += len(req.EmojiUpdates)
	totalSteps += len(req.PositionUpdates)
	if req.UpdatedTitle != nil {
		totalSteps++
	}
	if req.UpdatedIsPublic != nil {
		totalSteps++
	}
	return totalSteps
}

func editDeleteStage(deletedIDs []string, prog *editProgress) error {
	deletedCount := len(deletedIDs)
	for i, deleted := range deletedIDs {
		if err := telegram.DeleteSticker(deleted); err != nil {
			return err
		}
		prog.update(fmt.Sprintf("Deleting stickers (%d/%d)", i+1, deletedCount))
	}
	return nil
}

func editProcessStage(
	ctx context.Context,
	addedStickers []emote.EmoteInput,
	prog *editProgress,
) ([]telegram.InputSticker, error) {
	if len(addedStickers) == 0 {
		return nil, nil
	}

	prog.setMessage("Processing emotes")

	stickers, err := emotesToStickers(
		ctx,
		addedStickers,
		2,
		func(done, total int) {
			prog.current = prog.total - (total * 2) + done
			prog.notify(
				prog.current,
				prog.total,
				fmt.Sprintf("Processing emotes (%d/%d)", done, total),
			)
		},
	)
	if err != nil {
		return nil, err
	}

	return stickers, nil
}

func editUpdateIsPublicStage(
	req *EditPackRequest,
	packName string,
	prog *editProgress,
) error {
	if req.UpdatedIsPublic != nil {
		dbConn := config.Load().DBConn()
		err := dbConn.UpdateIsPublic(packName, *req.UpdatedIsPublic)
		if err != nil {
			return err
		}
		prog.update("Updated pack publicity")
	}
	return nil
}

func editUpdateTitleStage(
	req *EditPackRequest,
	pack *telegram.StickerPack,
	prog *editProgress,
) error {
	if req.UpdatedTitle != nil {
		if err := pack.SetTitle(*req.UpdatedTitle); err != nil {
			return err
		}
		prog.update("Updated pack title")
	}
	return nil
}

func editAddStage(
	ctx context.Context,
	pack *telegram.StickerPack,
	addedStickers []emote.EmoteInput,
	prog *editProgress,
) error {
	stickers, err := editProcessStage(ctx, addedStickers, prog)
	if err != nil {
		return fmt.Errorf("failed to process emotes: %w", err)
	}
	for i, sticker := range stickers {
		if err := pack.AddSticker(sticker); err != nil {
			return err
		}
		prog.update(fmt.Sprintf("Adding emotes (%d/%d)", i+1, len(stickers)))
	}
	return nil
}

func editEmojiStage(updates []StickerEmojiUpdate, prog *editProgress) error {
	updateCount := len(updates)
	for i, update := range updates {
		err := telegram.SetStickerEmojis(update.ID, update.Emojis)
		if err != nil {
			return err
		}
		prog.update(fmt.Sprintf("Updating emojis (%d/%d)", i+1, updateCount))
	}
	return nil
}

func editPositionStage(
	updates []StickerPositionUpdate,
	prog *editProgress,
) error {
	updateCount := len(updates)
	for i, update := range updates {
		err := telegram.SetStickerPosition(update.ID, update.Position)
		if err != nil {
			return err
		}
		prog.update(fmt.Sprintf("Updating positions (%d/%d)", i+1, updateCount))
	}
	return nil
}
