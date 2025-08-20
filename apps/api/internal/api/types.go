package api

import "github.com/Traunin/stickerpack-editor/apps/api/internal/db"

type GetPacksResponse struct {
	Packs []db.PackResponse `json:"packs"`
	Total int             `json:"total"`
}