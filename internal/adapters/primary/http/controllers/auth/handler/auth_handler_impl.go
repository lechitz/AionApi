// Package handler auth provides authentication handlers for login and logout functionalities.
package handler

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// Handler provides authentication controllers for login and logout functionalities.
type Handler struct {
	Service input.AuthService
	Logger  output.ContextLogger
	Config  *config.Config
}

// New initializes and returns a new Auth instance with AuthService and Logger dependencies.
func New(authService input.AuthService, cfg *config.Config, logger output.ContextLogger) *Handler {
	return &Handler{
		Service: authService,
		Config:  cfg,
		Logger:  logger,
	}
}
