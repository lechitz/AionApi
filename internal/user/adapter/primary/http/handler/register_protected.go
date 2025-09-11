package handler

import (
	"net/http"

	authMiddleware "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	authinput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTPProtected registers PROTECTED endpoints for the user context.
// Applies AuthMiddleware only to the subtree that requires credentials.
func RegisterHTTPProtected(r ports.Router, h *Handler, authService authinput.AuthService, lg logger.ContextLogger) {
	r.Group("/v1/users", func(ur ports.Router) {
		mw := authMiddleware.New(authService, lg)

		ur.GroupWith(mw.Auth, func(pr ports.Router) {
			// Protected reads (consider RBAC)
			pr.GET("/all", http.HandlerFunc(h.ListAll))           // Recommended: admin-only
			pr.GET("/{user_id}", http.HandlerFunc(h.GetUserByID)) // Recommended: self or admin

			// Protectedly writes
			pr.PUT("/", http.HandlerFunc(h.UpdateUser))
			pr.PUT("/password", http.HandlerFunc(h.UpdateUserPassword))
			pr.DELETE("/", http.HandlerFunc(h.SoftDeleteUser))
		})
	})
}
