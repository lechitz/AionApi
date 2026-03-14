package handler

import (
	"net/http"

	authmw "github.com/lechitz/AionApi/internal/auth/adapter/primary/http/middleware"
	authInput "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// RegisterHTTP mounts realtime routes behind the auth middleware.
func RegisterHTTP(r ports.Router, h *Handler, authService authInput.AuthService, log logger.ContextLogger) {
	if h == nil || authService == nil {
		return
	}

	r.Group("/realtime", func(rr ports.Router) {
		rr.Use(authmw.New(authService, log).Auth)
		rr.GET(h.streamPath(), http.HandlerFunc(h.Stream))
	})
}
