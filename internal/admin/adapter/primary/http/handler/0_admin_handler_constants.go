// Package handler is the handler for the admin context in the application.
package handler

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the tracer name for admin handler.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.admin.handler"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanHealthCheck is the span name for health check handler.
	SpanHealthCheck = "admin.handler.health_check"

	// SpanGetMetrics is the span name for metrics handler.
	SpanGetMetrics = "admin.handler.get_metrics"

	// SpanGetInfo is the span name for info handler.
	SpanGetInfo = "admin.handler.get_info"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusHealthy indicates the system is healthy.
	StatusHealthy = "healthy"

	// StatusUnhealthy indicates the system is unhealthy.
	StatusUnhealthy = "unhealthy"
)

// =============================================================================
// BUSINESS LOGIC - Messages
// =============================================================================

const (
	// MsgHealthCheckSuccess is the message for successful health check.
	MsgHealthCheckSuccess = "health check passed"

	// MsgHealthCheckFailed is the message for failed health check.
	MsgHealthCheckFailed = "health check failed"
)
