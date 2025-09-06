package handler

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
)

type Handler struct {
	TagService input.TagService
	Logger     logger.ContextLogger
}

func NewHandler(svc input.TagService, logger logger.ContextLogger) *Handler {
	return &Handler{TagService: svc, Logger: logger}
}
