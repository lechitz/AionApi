// Package handler auth provides authentication handlers for login and logout functionalities.
package handler

import (
	"github.com/lechitz/AionApi/internal/auth/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Handler provides authentication controllers for login and logout functionalities.
type Handler struct {
	Service input.AuthService
	Logger  logger.ContextLogger
	Config  *config.Config
}

// New initializes and returns a new Auth instance with AuthService and Logger dependencies.
func New(authService input.AuthService, cfg *config.Config, logger logger.ContextLogger) *Handler {
	return &Handler{
		Service: authService,
		Config:  cfg,
		Logger:  logger,
	}
}
