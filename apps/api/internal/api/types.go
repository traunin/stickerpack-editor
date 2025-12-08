package api

import (
	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
)

type GetPacksResponse struct {
	Packs []telegram.PackPreview `json:"packs"`
	Total int               `json:"total"`
}
