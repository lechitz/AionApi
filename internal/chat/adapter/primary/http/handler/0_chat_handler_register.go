// Package handler registers HTTP routes for the chat module.
package handler

import (
	"net/http"

	authMiddleware "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	"github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTP attaches the chat routes to the given router.
func RegisterHTTP(r ports.Router, h *Handler, authService input.AuthService, log logger.ContextLogger) {

	if authService != nil {
		mw := authMiddleware.New(authService, log)
		r.GroupWith(mw.Auth, func(pr ports.Router) {
			pr.POST("/chat", http.HandlerFunc(h.Chat))
		})
	}
}
