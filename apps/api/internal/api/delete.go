package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/telegram"
)

type DeletePackRequest struct {
	UserID   int64
	PackName string `json:"pack_name"`
}

type DeletePackResponse struct {
	Success bool `json:"success"`
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
