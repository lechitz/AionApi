// Package handler is the handler for the auth context in the application.
package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTP registers the auth context's HTTP handlers with the provided router.
func RegisterHTTP(r ports.Router, h *Handler) {
	r.Group("/auth", func(ar ports.Router) {
		// Public.
		ar.POST("/login", http.HandlerFunc(h.Login))

		// Protected.
		authmw := middleware.New(h.Service, h.Logger)
		ar.GroupWith(authmw.Auth, func(pr ports.Router) {
			pr.POST("/logout", http.HandlerFunc(h.Logout))
		})
	})
}
