package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// SetupCORS creates and configures CORS middleware
func SetupCORS(allowedOrigins []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})
}
