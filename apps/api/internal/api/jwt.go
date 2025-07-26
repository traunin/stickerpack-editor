package api

import (
	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/golang-jwt/jwt"
)

func sign(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	key := []byte(config.Load().SecretKey)

	return token.SignedString(key)
}
