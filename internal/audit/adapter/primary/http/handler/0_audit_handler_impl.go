package handler

import (
	"github.com/lechitz/aion-api/internal/audit/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// Handler wires audit use cases to HTTP handlers.
type Handler struct {
	Service input.Service
	Logger  logger.ContextLogger
	Config  *config.Config
}

// New creates a new audit HTTP handler.
func New(service input.Service, cfg *config.Config, log logger.ContextLogger) *Handler {
	return &Handler{
		Service: service,
		Config:  cfg,
		Logger:  log,
	}
}
