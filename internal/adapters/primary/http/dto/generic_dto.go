// Package dto contains Data Transfer Objects used by the HTTP layer.
package dto

import (
	"time"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
)

// HealthCheckResponse represents the health check response structure.
type HealthCheckResponse struct {
	Name      string                 `json:"name"`
	Env       string                 `json:"env"`
	Version   string                 `json:"version"`
	Timestamp time.Time              `json:"timestamp"`
	Status    constants.HealthStatus `json:"status"`
}

// ErrorResponse represents the error response structure.
type ErrorResponse struct {
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Status    constants.HealthStatus `json:"status"`
}
