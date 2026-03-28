// Package handler registers HTTP routes for the admin context.
package handler

import (
	"net/http"

	authMiddleware "github.com/lechitz/aion-api/internal/auth/adapter/primary/http/middleware"
	authinput "github.com/lechitz/aion-api/internal/auth/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/platform/server/http/ports"
)

// RegisterHTTP registers the admin-related HTTP handlers with the provided router.
// All admin endpoints require authentication and admin/owner role.
func RegisterHTTP(r ports.Router, h *Handler, authService authinput.AuthService, lg logger.ContextLogger) {
	if authService != nil {
		// All admin routes require authentication
		mw := authMiddleware.New(authService, lg)
		r.GroupWith(mw.Auth, func(ar ports.Router) {
			// Admin routes under /admin/users/{user_id}
			ar.PUT("/admin/users/{user_id}/roles", http.HandlerFunc(h.UpdateUserRoles))

			// Role management with hierarchy validation
			ar.PUT("/admin/users/{user_id}/promote-admin", http.HandlerFunc(h.PromoteToAdmin))
			ar.PUT("/admin/users/{user_id}/demote-admin", http.HandlerFunc(h.DemoteFromAdmin))
			ar.PUT("/admin/users/{user_id}/block", http.HandlerFunc(h.BlockUser))
			ar.PUT("/admin/users/{user_id}/unblock", http.HandlerFunc(h.UnblockUser))
		})
	}
}
