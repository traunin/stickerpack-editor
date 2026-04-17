package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const UserIDContextKey = "userID"

type NoAuthRoute struct {
	Path        string
	Method      string
	PrefixMatch bool
}

type Claims struct {
	Id int64 `json:"id"`
}

func SignID(id int64, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})

	return token.SignedString(key)
}

func DecodeID(tokenStr string, key []byte) (int64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}
			return key, nil
		})

	if err != nil {
		return 0, fmt.Errorf("error parsing JWT: %w", err)
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("no id in JWT")
	}

	return int64(id), nil
}

func isNoAuth(noAuthRoutes []NoAuthRoute, r *http.Request) bool {
	for _, route := range noAuthRoutes {
		if route.Method != "" && r.Method != route.Method {
			continue
		}
		if route.PrefixMatch && strings.HasPrefix(r.URL.Path, route.Path) {
			return true
		}
		if !route.PrefixMatch && r.URL.Path == route.Path {
			return true
		}
	}
	return false
}

func (h *Handler) jwtMiddleware(noAuthRoutes []NoAuthRoute, next http.Handler) http.Handler {
	key := []byte(h.cfg.SecretKey())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isNoAuth(noAuthRoutes, r) {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("jwt")
		if err != nil || cookie.Value == "" {
			http.Error(w, "JWT cookie required", http.StatusUnauthorized)
			return
		}

		userID, err := DecodeID(cookie.Value, key)
		if err != nil {
			errorMessage := fmt.Sprintf("Invalid JWT token: %v", err)
			http.Error(w, errorMessage, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(r *http.Request) (int64, error) {
	userID, ok := r.Context().Value(UserIDContextKey).(int64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
