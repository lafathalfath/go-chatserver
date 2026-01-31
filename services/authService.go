package services

import (
	"context"
	"errors"
	"github.com/lafathalfath/go-chatserver/cache"
	contextkeys "github.com/lafathalfath/go-chatserver/context-keys"
	"github.com/lafathalfath/go-chatserver/database"
	"github.com/lafathalfath/go-chatserver/graph/models"
	"github.com/lafathalfath/go-chatserver/helpers"
	"github.com/lafathalfath/go-chatserver/validators"
	"net/http"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, input *models.Login) (bool, error)
	Refresh(ctx context.Context) (bool, error)
	Me(ctx context.Context) (*models.User, error)
	Logout(ctx context.Context) (bool, error)
}

type authService struct {
	*database.ConnectionObj
}

func AuthServices() AuthService {
	conn := &database.DBConnection
	return &authService{conn}
}

var secureCookie bool = helpers.Env("ENV") == "development"

func (s *authService) Login(ctx context.Context, input *models.Login) (bool, error) {
	user, err := validators.ValidateLogin(s.DB, *input)
	if err != nil {
		return false, err
	}

	accessToken, _ := helpers.GenerateAccessToken(user.ID, time.Hour)
	refreshToken, _ := helpers.GenerateRefreshToken(user.ID, 24*time.Hour)

	w := ctx.Value(contextkeys.ContextKeyWriter).(http.ResponseWriter)
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   secureCookie,
		// SameSite: http.SameSiteLaxMode,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   secureCookie,
		// SameSite: http.SameSiteLaxMode,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
	})
	return true, nil
}

func (s *authService) Refresh(ctx context.Context) (bool, error) {
	r := ctx.Value(contextkeys.ContextKeyRequest).(*http.Request)
	w := ctx.Value(contextkeys.ContextKeyWriter).(http.ResponseWriter)

	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		return false, errors.New("Unauthorized")
	}
	userID, err := helpers.ParseRefreshToken(cookie.Value)
	if err != nil {
		return false, errors.New("Unauthorized")
	}
	newAccessToken, _ := helpers.GenerateAccessToken(userID, time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    newAccessToken,
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   60 * 60,
	})
	return true, nil
}

func (s *authService) Me(ctx context.Context) (*models.User, error) {
	userId, ok := helpers.GetUserId(ctx)
	if !ok {
		return nil, errors.New("Unauthorized")
	}
	cacheKey := "users:" + userId
	var user models.User
	if err := cache.Scan(cacheKey, &user); err == nil {
		return &user, nil
	} else if err := s.DB.Select("id", "name", "email", "password").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}
	cache.Set(cacheKey, user, 10*time.Minute)
	return &user, nil
}

func (s *authService) Logout(ctx context.Context) (bool, error) {
	w := ctx.Value(contextkeys.ContextKeyWriter).(http.ResponseWriter)
	clear := func(name, path string) {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			MaxAge:   -1,
			HttpOnly: true,
		})
	}
	clear("accessToken", "/")
	clear("refreshToken", "/")
	return true, nil
}
