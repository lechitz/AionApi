package handler

import (
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
)

// Handler is the handler for user-related HTTP operations.
type Handler struct {
	UserService input.UserService
	Logger      logger.ContextLogger
	Config      *config.Config
}

// New returns a User handler with dependencies injected.
func New(userService input.UserService, cfg *config.Config, logger logger.ContextLogger) *Handler {
	return &Handler{
		UserService: userService,
		Config:      cfg,
		Logger:      logger,
	}
}
