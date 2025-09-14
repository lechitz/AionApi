package handler

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
)

// Handler is the controller for the tag service.
type Handler struct {
	TagService input.TagService
	Logger     logger.ContextLogger
}

// NewHandler wires dependencies and returns a Controller.
func NewHandler(svc input.TagService, logger logger.ContextLogger) *Handler {
	return &Handler{TagService: svc, Logger: logger}
}
