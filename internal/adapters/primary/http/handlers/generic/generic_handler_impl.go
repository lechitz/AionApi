// Package generic provides common HTTP handlers for the application.
package generic

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// Handler provides common HTTP handlers for the application.
type Handler struct {
	Logger        output.ContextLogger
	GeneralConfig config.GeneralConfig
}

// New initializes and returns a new Generic instance with a Logger dependency.
func New(logger output.ContextLogger, general config.GeneralConfig) *Handler {
	return &Handler{
		Logger:        logger,
		GeneralConfig: general,
	}
}
