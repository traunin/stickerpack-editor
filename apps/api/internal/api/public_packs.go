package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
)

func getPublicPacksHandler(w http.ResponseWriter, r *http.Request) {
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

	resp, err := publicPacksPreviews(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func publicPacksPreviews(
	ctx context.Context,
	page,
	pageSize int,
) (*GetPacksResponse, error) {
	db := config.Load().DBConn()
	packs, err := db.PublicStickerpacks(page, pageSize)
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

	total, err := db.PublicPacksCount()
	if err != nil {
		return nil, err
	}

	return &GetPacksResponse{
		Packs: previews,
		Total: total,
	}, nil
}
