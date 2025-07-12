// Package commonkeys contains constants used throughout the application.
package commonkeys

// --- Application Metadata.
const (
	APIName    = "api_name"    // Application or service name
	AppEnv     = "app_env"     // Application environment ("development", "production", etc)
	AppVersion = "app_version" // Application version
	Setting    = "setting"     // Generic configuration setting
)

// --- HTTP Server/Request.
const (
	ServerHTTPName = "server_name" // HTTP server name
	ServerHTTPAddr = "http_port"   // HTTP server port/address
	URLPath        = "path"        // Current URL path
	Method         = "method"      // HTTP method (GET, POST, etc)
	Status         = "status"      // Status of operation/process (e.g., success, failure)
)

// --- Database.
const (
	DBName                     = "db_name" // Database name
	DBHost                     = "host"    // Database host
	DBPort                     = "port"    // Database port
	DBTryConnectingWithRetries = "try"     // Connection retry counter
)

// CacheAddr --- Cache.
const (
	CacheAddr = "cache_addr" // Cache/Redis server address.
)

// --- Request/Tracing Identifiers.
const (
	RequestID  = "request_id"   // Internal/external request ID in logs, headers, or context
	XRequestID = "X-Request-ID" // Header key for external request tracking
)

// Input --- GraphQL.
const (
	Input = "input" // Input value (GraphQL, forms, payloads)
)

// Error --- General tag.
const (
	Error = "error" // Generic error, often used in logging or context
)
