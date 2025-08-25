// Package handler generic provides common HTTP controllers for the application.
package handler

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// Handler provides common HTTP controllers for the application.
type Handler struct {
	Logger        output.ContextLogger
	GeneralConfig config.GeneralConfig
}

// New initializes and returns a new Generic instance with a Logger dependency.
func New(logger output.ContextLogger, generalCfg config.GeneralConfig) *Handler {
	return &Handler{
		Logger:        logger,
		GeneralConfig: generalCfg,
	}
}
