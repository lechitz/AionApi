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

// Server keys.
const (
	// ServerName is the HTTP server name.
	ServerName = "server_name"

	// ServerAddr is the HTTP server ports/address.
	ServerAddr = "http_port"
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

	// OperationCreate is the name of the create operation.
	OperationCreate = "create"

	// OperationSave is the name of the save operation.
	OperationSave = "save"

	// OperationGet is the name of the get operation.
	OperationGet = "get"

	// OperationUpdate is the name of the update operation.
	OperationUpdate = "update"

	// OperationDelete is the name of the delete operation.
	OperationDelete = "delete"

	// Entity is the name of the entity being operated on.
	Entity = "entity"

	EntityToken = "token"

	// StatusSuccess is the standard value for successful status.
	StatusSuccess = "success"

	// Error is the generic error key for logging/context.
	Error = "error"

	// Fields is used in controllers for payload fields.
	Fields = "fields"
)

// Database keys.
const (
	// DBName is the database name.
	DBName = "db_name"

	// DBHost is the database host.
	DBHost = "host"

	// DBPort is the database ports.
	DBPort = "ports"

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
