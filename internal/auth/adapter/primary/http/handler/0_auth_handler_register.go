// Package handler registers HTTP routes for the authentication context.
package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTP attaches the auth routes to the given router.
func RegisterHTTP(r ports.Router, h *Handler) {
	r.Group("/auth", func(ar ports.Router) {
		// Public endpoint.
		ar.POST("/login", http.HandlerFunc(h.Login))

		// Protected endpoints require an authenticated context.
		authmw := middleware.New(h.Service, h.Logger)
		ar.GroupWith(authmw.Auth, func(pr ports.Router) {
			pr.POST("/logout", http.HandlerFunc(h.Logout))
		})
	})
}
