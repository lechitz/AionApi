// Package handlers provide common HTTP handlers for the application.
package handlers

import (
	"net/http"

	"github.com/lechitz/AionApi/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"

	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// Generic represents a type that provides common handlers for HTTP requests with logging capabilities.
type Generic struct {
	Logger logger.Logger
}

// NewGeneric initializes and returns a new Generic instance with the provided logger dependency.
func NewGeneric(logger logger.Logger) *Generic {
	return &Generic{Logger: logger}
}

// HealthCheckHandler handles HTTP requests for checking the health of the service and responds with a healthy status message.
func (h *Generic) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	response.Return(w, http.StatusOK, []byte(constants.MsgServiceIsHealthy), h.Logger)
}

// NotFoundHandler handles HTTP requests for non-existent resources and responds with a 404 status and a resource not found message.
func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Return(w, http.StatusNotFound, []byte(constants.MsgResourceNotFound), h.Logger)
	h.Logger.Infow(constants.MsgResourceNotFound, "path", r.URL.Path)
}
