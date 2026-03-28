package handler

import (
	"net/http"

	authMiddleware "github.com/lechitz/aion-api/internal/auth/adapter/primary/http/middleware"
	authinput "github.com/lechitz/aion-api/internal/auth/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/platform/server/http/ports"
)

// RegisterHTTP registers audit routes with auth protection.
func RegisterHTTP(r ports.Router, h *Handler, authService authinput.AuthService, lg logger.ContextLogger) {
	if authService == nil {
		return
	}

	mw := authMiddleware.New(authService, lg)
	r.GroupWith(mw.Auth, func(ar ports.Router) {
		ar.GET("/audit/events", http.HandlerFunc(h.ListEvents))
	})
}
