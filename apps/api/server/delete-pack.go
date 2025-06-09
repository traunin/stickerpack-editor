package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/telegram"
)

type DeletePackRequest struct {
	PackName string `json:"pack_name"`
	// UserID   string `json:"userID"`
}

type DeletePackResponse struct {
	Success bool `json:"success"`
}

func deletePackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method is not DELETE", http.StatusBadRequest)
		return
	}

	var req DeletePackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON schema", http.StatusBadRequest)
		return
	}

	if req.PackName == "" {
		http.Error(w, "stickerpack name missing", http.StatusBadRequest)
		return
	}

	pack := telegram.StickerPack{
		Name: req.PackName,
	}

	err := pack.Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete sticker pack: %v", err), http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode(DeletePackResponse{Success: true})
}
