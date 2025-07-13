// Package commonkeys contains constants used throughout the application.
package commonkeys

// Application metadata keys.
const (
	// APIName is the application or service name.
	APIName = "api_name"

	// AppEnv is the application environment ("development", "production", etc).
	AppEnv = "app_env"

	// AppVersion is the application version.
	AppVersion = "app_version"

	// Setting is a generic configuration setting.
	Setting = "setting"
)

// HTTP server keys.
const (
	// ServerHTTPName is the HTTP server name.
	ServerHTTPName = "server_name"

	// ServerHTTPAddr is the HTTP server port/address.
	ServerHTTPAddr = "http_port"
)

// HTTP request keys.
const (
	// URLPath is the current URL path.
	URLPath = "path"

	// Method is the HTTP method (GET, POST, etc).
	Method = "method"
)

// Operation/status keys.
const (
	// Status is the status of an operation / process (e.g., success, failure).
	Status = "status"

	// Operation is the name of operation/handler/action.
	Operation = "operation"

	// StatusSuccess is the standard value for successful status.
	StatusSuccess = "success"

	// Error is the generic error key for logging/context.
	Error = "error"

	// Fields is used in handlers for payload fields.
	Fields = "fields"
)

// Database keys.
const (
	// DBName is the database name.
	DBName = "db_name"

	// DBHost is the database host.
	DBHost = "host"

	// DBPort is the database port.
	DBPort = "port"

	// DBTryConnectingWithRetries is the connection retry counter.
	DBTryConnectingWithRetries = "try"
)

// Cache keys.
const (
	// CacheAddr is the cache/Redis server address.
	CacheAddr = "cache_addr"
)

// Request/tracing identifiers.
const (
	// RequestID is the internal/external request ID in logs, headers, or context.
	RequestID = "request_id"

	// XRequestID is the header key for external request tracking.
	XRequestID = "X-Request-ID"
)

// Input keys (GraphQL, forms, payloads).
const (
	// Input is the key for input value (GraphQL, forms, payloads).
	Input = "input"
)
