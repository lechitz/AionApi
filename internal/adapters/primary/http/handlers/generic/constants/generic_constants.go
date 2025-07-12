// Package constants contains constants used throughout the generic handler.
package constants

import "errors"

// Tracer names for OpenTelemetry generic handler operations.
const (
	TracerGenericHandler     = "aionapi.generic.handler"  // Main tracer for generic handler
	TracerHealthCheckHandler = "generic.health_check"     // Span name for health check
	TracerErrorHandler       = "generic.error_handler"    // Span name for internal error handler
	TracerRecoveryHandler    = "generic.recovery_handler" // Span name for recovery from panic
)

// EventHealthCheck Event names for key points within generic handler spans.
const (
	EventHealthCheck = "health_check_event"
)

// StatusHealthCheckOK Status names for semantic span states.
const (
	StatusHealthCheckOK = "health_check_ok" // Semantic status for a healthy state
)

// HealthStatusHealthy Other string constants for health status.
const (
	HealthStatusHealthy = "healthy"
)

// Error variables for generic handler (to be used as an error interface).
var (
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrResourceNotFound   = errors.New("resource not found")
	ErrRecoveredFromPanic = errors.New("recovered from panic")
	ErrInternalServer     = errors.New("internal server error")
)

// Standardized messages for logs, responses, and traces.
const (
	MsgServiceIsHealthy    = "service is healthy"
	MsgMethodNotAllowed    = "the requested method is not allowed"
	MsgResourceNotFound    = "resource not found"
	MsgRecoveredFromPanic  = "application recovered from panic"
	MsgInternalServerError = "internal server error"
)

// StacktraceFormat is the printf format for recording recovered panic and stacktrace in tracing/logs.
const StacktraceFormat = "%v\n%s"
