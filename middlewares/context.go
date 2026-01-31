package middlewares

import (
	"context"
	contextkeys "github.com/lafathalfath/go-chatserver/context-keys"
	"net/http"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextkeys.ContextKeyWriter, w)

		ctx = context.WithValue(ctx, contextkeys.ContextKeyRequest, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
