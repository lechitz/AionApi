package handler

import (
	"net/http"

	authMiddleware "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	authinput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTP registers the user-related HTTP handlers with the provided router.
func RegisterHTTP(r ports.Router, h *Handler, authService authinput.AuthService, lg logger.ContextLogger) {
	r.Group("/user", func(ur ports.Router) {
		// Public
		ur.POST("/create", http.HandlerFunc(h.Create))

		// Private
		if authService != nil {
			mw := authMiddleware.New(authService, lg)
			ur.GroupWith(mw.Auth, func(pr ports.Router) {
				pr.GET("/all", http.HandlerFunc(h.ListAll))
				pr.GET("/{user_id}", http.HandlerFunc(h.GetUserByID))
				pr.PUT("/", http.HandlerFunc(h.UpdateUser))
				pr.PUT("/password", http.HandlerFunc(h.UpdateUserPassword))
				pr.DELETE("/", http.HandlerFunc(h.SoftDeleteUser))
			})
		}
	})
}
