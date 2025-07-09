package constants

// HealthStatus represents possible health statuses for the service.
type HealthStatus string

const (
	// HealthStatusHealthy indicates the service is operational.
	HealthStatusHealthy HealthStatus = "healthy"

	// HealthStatusUnhealthy indicates the service is unhealthy.
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

const (

	// MsgSpanHealthCheckOK describes a successful health check span.
	MsgSpanHealthCheckOK = "health check successful"

	// MsgServiceIsHealthy is the message returned when the service is healthy.
	MsgServiceIsHealthy = "Service is healthy"

	// MsgResourceNotFound is the message returned when a resource is not found.
	MsgResourceNotFound = "Resource not found"
)

const (
	// TracerGenericHandler is the name of the tracer for the generic handler.
	TracerGenericHandler = "GenericHandler"

	// TracerHealthCheckHandler is the name of the tracer for the health check handler.
	TracerHealthCheckHandler = "HealthCheckHandler"
)
