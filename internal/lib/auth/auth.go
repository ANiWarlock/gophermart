package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/ANiWarlock/gophermart/cmd/gophermart/config"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
)

var ErrInvalidToken = errors.New("token is not valid")

var secret string

type ctxKey uint

const CtxKeyUserID ctxKey = iota

func SetSecretKey(cfg *config.AppConfig) {
	secret = cfg.SecretKey
	if secret == "" {
		secret = "mySecretKey"
	}
}

func GetUserID(tokenString string) (uint, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token string: %w", err)
	}
	if !token.Valid {
		log.Printf("Received invalid token: %s \n", tokenString)
		return 0, ErrInvalidToken
	}

	userID := uint(claims["userId"].(float64))

	return userID, nil
}

func BuildCookieStringValue(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userID,
		})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed signing string: %w", err)
	}

	return tokenString, nil
}

func CurrentUser(ctx context.Context) uint {
	if userID := ctx.Value(CtxKeyUserID); userID != nil {
		return userID.(uint)
	}

	return 0
}

func SetCookie(userID uint, rw http.ResponseWriter) error {
	cookieStringValue, err := BuildCookieStringValue(userID)
	if err != nil {
		return fmt.Errorf("failed to build cookie string: %w", err)
	}

	newCookie := &http.Cookie{
		Name:     "auth",
		Value:    cookieStringValue,
		HttpOnly: true,
	}

	http.SetCookie(rw, newCookie)
	return nil
}
