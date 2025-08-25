package handler

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// Handler is the handler for user-related HTTP operations.
type Handler struct {
	UserService input.UserService
	Logger      output.ContextLogger
	Config      *config.Config
}

// New returns a User handler with dependencies injected.
func New(userService input.UserService, cfg *config.Config, logger output.ContextLogger) *Handler {
	return &Handler{
		UserService: userService,
		Config:      cfg,
		Logger:      logger,
	}
}
