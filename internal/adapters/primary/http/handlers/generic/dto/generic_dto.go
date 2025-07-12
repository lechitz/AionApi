// Package dto contains Data Transfer Objects used by the HTTP layer.
package dto

import (
	"time"
)

// HealthCheckResponse represents the health check response structure.
type HealthCheckResponse struct {
	Name      string    `json:"name"`
	Env       string    `json:"env"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}
