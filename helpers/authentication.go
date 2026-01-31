package helpers

import (
	"context"
	"errors"
	contextkeys "github.com/lafathalfath/go-chatserver/context-keys"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	secretKey  = []byte(Env("SECRET_KEY"))
	refreshKey = []byte(Env("REFRESH_KEY"))
)

func GenerateAccessToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshKey)
}

func ParseAccessToken(jwtToken string) (string, error) {
	return parseToken(jwtToken, secretKey)
}

func ParseRefreshToken(jwtToken string) (string, error) {
	return parseToken(jwtToken, refreshKey)
}

func parseToken(jwtToken string, key []byte) (string, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return key, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("Invalid Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid claims")
	}
	userId, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("User ID not found")
	}
	return userId, nil
}

func GetUserId(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(contextkeys.UserIDContextKey).(string)
	return userID, ok
}
