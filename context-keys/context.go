package contextkeys

type contextKey string

var (
	ContextKeyRequest contextKey = "request"
	ContextKeyWriter contextKey = "writer"
	UserIDContextKey contextKey = "UserID"
)