package tags

import (
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

type Handler struct {
	TagService input.TagService
	Logger     output.ContextLogger
}

func NewHandler(svc input.TagService, logger output.ContextLogger) *Handler {
	return &Handler{TagService: svc, Logger: logger}
}
