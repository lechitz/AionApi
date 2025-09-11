package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTPPublic registers PUBLIC endpoints for the user context.
// The Platform composer will mount this group under the global API prefix.
func RegisterHTTPPublic(r ports.Router, h *Handler) {
	r.Group("/v1/users", func(ur ports.Router) {
		// Public (if your product allows open sign-up)
		ur.POST("/create", http.HandlerFunc(h.Create))
	})
}
