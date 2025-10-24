// Package cors provides HTTP middleware for CORS configuration.
package cors

import (
	"net/http"

	"github.com/go-chi/cors"
)

// New returns a CORS middleware handler with the recommended options for the AionApi frontend integration.
func New() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
