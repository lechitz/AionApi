// Package handler exposes realtime HTTP endpoints.
package handler

import (
	"time"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	realtimeInput "github.com/lechitz/aion-api/internal/realtime/core/ports/input"
)

// Handler serves realtime SSE endpoints.
type Handler struct {
	Service realtimeInput.Service
	Config  *config.Config
	Logger  logger.ContextLogger
}

// New creates a realtime HTTP handler.
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
