package handler

import (
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	realtimeInput "github.com/lechitz/AionApi/internal/realtime/core/ports/input"
)

type Handler struct {
	Service realtimeInput.Service
	Config  *config.Config
	Logger  logger.ContextLogger
}

func New(service realtimeInput.Service, cfg *config.Config, log logger.ContextLogger) *Handler {
	return &Handler{
		Service: service,
		Config:  cfg,
		Logger:  log,
	}
}

func (h *Handler) heartbeatInterval() time.Duration {
	if h.Config == nil || h.Config.Realtime.HeartbeatInterval <= 0 {
		return 15 * time.Second
	}
	return h.Config.Realtime.HeartbeatInterval
}

func (h *Handler) streamPath() string {
	if h.Config == nil || h.Config.Realtime.StreamPath == "" {
		return streamRoute
	}
	return h.Config.Realtime.StreamPath
}
