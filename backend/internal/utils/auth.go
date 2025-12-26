package utils

import (
	"context"
	"net/http"

	"github.com/Ankitchan/hackernews-clone/pkg/auth"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserContextKey is the key for storing user claims in context
	UserContextKey ContextKey = "user"
)

// GetUserFromContext extracts user claims from request context
func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*auth.Claims)
	return claims, ok
}

// GetUserFromRequest extracts user claims from HTTP request context
func GetUserFromRequest(r *http.Request) (*auth.Claims, bool) {
	return GetUserFromContext(r.Context())
}
