package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ankitchan/hackernews-clone/internal/utils"
	"github.com/Ankitchan/hackernews-clone/pkg/auth"
)

// AuthResponse represents error response structure
type AuthResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// AuthMiddleware validates JWT tokens and adds user info to request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Add user claims to request context
		ctx := context.WithValue(r.Context(), utils.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware validates JWT tokens if present, but allows requests without auth
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No auth header, continue without user context
			next.ServeHTTP(w, r)
			return
		}

		// Try to parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]

			// Try to validate token
			if claims, err := auth.ValidateToken(tokenString); err == nil {
				// Valid token, add to context
				ctx := context.WithValue(r.Context(), utils.UserContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Invalid or missing token, continue without user context
		next.ServeHTTP(w, r)
	})
}


// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(AuthResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}
