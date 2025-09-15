// Package handler provides HTTP handlers for the authentication context.
package handler

import (
	"github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// @tag.name Auth
// @tag.description Authentication endpoints (login/logout). Supports Bearer token and session cookie.

// Handler wires authentication use cases to HTTP handlers.
type Handler struct {
	Service input.AuthService
	Logger  logger.ContextLogger
	Config  *config.Config
}

// New creates a new Handler instance with its dependencies wired.
// Using 'log' as param name to avoid shadowing the imported 'logger' package.
func New(authService input.AuthService, cfg *config.Config, log logger.ContextLogger) *Handler {
	return &Handler{
		Service: authService,
		Config:  cfg,
		Logger:  log,
	}
}
