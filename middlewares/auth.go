package middlewares

import (
	"context"
	contextkeys "github.com/lafathalfath/go-chatserver/context-keys"
	"github.com/lafathalfath/go-chatserver/helpers"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accessToken")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		token := cookie.Value
		userID, err := helpers.ParseAccessToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), contextkeys.UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
