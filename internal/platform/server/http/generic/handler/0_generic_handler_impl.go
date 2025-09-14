// Package handler generic provides common HTTP controllers for the application.
package handler

import (
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Handler provides common HTTP controllers for the application.
type Handler struct {
	Logger        logger.ContextLogger
	GeneralConfig config.GeneralConfig
}

// New initializes and returns a new Generic instance with a Logger dependency.
func New(logger logger.ContextLogger, generalCfg config.GeneralConfig) *Handler {
	return &Handler{
		Logger:        logger,
		GeneralConfig: generalCfg,
	}
}
