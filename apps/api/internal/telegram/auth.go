package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
)

type AuthRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

func ParseAuth(r *http.Request) (*AuthRequest, error) {
	req := new(AuthRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	if err := validateRequest(req); err != nil {
		return nil, err
	}

	return req, nil
}

func validateRequest(req *AuthRequest) error {
	// https://core.telegram.org/widgets/login#checking-authorization
	dataCheckString := buildCheckString(req)

	secret := sha256.Sum256([]byte(config.Load().TelegramToken()))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))

	botMAC := mac.Sum(nil)
	authMAC, err := hex.DecodeString(req.Hash)
	if err != nil {
		return fmt.Errorf("invalid hash")
	}

	if hmac.Equal(botMAC, authMAC) {
		return nil
	}

	return fmt.Errorf("unable to verify auth data")
}

func buildCheckString(req *AuthRequest) string {
	fields := map[string]string{
		"auth_date":  strconv.FormatInt(req.AuthDate, 10),
		"first_name": req.FirstName,
		"id":         strconv.FormatInt(req.ID, 10),
		"last_name":  req.LastName,
		"photo_url":  req.PhotoURL,
		"username":   req.Username,
	}

	var keyvals []string
	for k, v := range fields {
		if v != "" {
			keyval := fmt.Sprintf("%s=%s", k, v)
			keyvals = append(keyvals, keyval)
		}
	}

	sort.Strings(keyvals)
	return strings.Join(keyvals, "\n")
}
