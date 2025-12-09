// Package usecase contains business logic for the Admin context.
package usecase

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer used for the admin use case.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.admin.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanHealthCheck is the span name for health check operations.
	SpanHealthCheck = "admin.health_check"

	// SpanGetMetrics is the span name for getting metrics.
	SpanGetMetrics = "admin.get_metrics"

	// SpanGetInfo is the span name for getting application info.
	SpanGetInfo = "admin.get_info"
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
