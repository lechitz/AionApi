// Package config constants contains constants used throughout the application.
package config

import "time"

// Validation constants.
const (

	// MinHTTPTimeout is the minimum allowed timeout for HTTP requests.
	MinHTTPTimeout = 1 * time.Second

	// MinGraphQLTimeout is the minimum allowed timeout for GraphQL requests.
	MinGraphQLTimeout = 1 * time.Second

	// MinContextRequest is the minimum allowed timeout for context requests.
	MinContextRequest = 500 * time.Millisecond

	// MinDBRetryInterval is the minimum allowed interval for database retries.
	MinDBRetryInterval = 1 * time.Second

	// MinShutdownTimeout is the minimum allowed timeout for shutdown operations.
	MinShutdownTimeout = 1 * time.Second

	// MinCachePoolSize is the minimum allowed size for the cache pool.
	MinCachePoolSize = 1

	// MinDBMaxOpenConns is the minimum allowed value for the maximum number of open connections.
	MinDBMaxOpenConns = 1

	// MinDBMaxIdleConns is the minimum allowed value for the maximum number of idle connections.
	MinDBMaxIdleConns = 0

	// MinDBConnMaxLifetimeMin is the minimum allowed value for the maximum connection lifetime.
	MinDBConnMaxLifetimeMin = 0

	// MinDBMaxRetries is the minimum allowed value for the maximum number of retries.
	MinDBMaxRetries = 1
)

// ErrFailedToProcessEnvVars is returned when environment variables cannot be processed.
const ErrFailedToProcessEnvVars = "failed to process environment variables: %v"

// ErrGenerateSecretKey is returned when the secret key generation fails.
const ErrGenerateSecretKey = "failed to generate secret key"

// SecretKeyWasNotSet is logged when no SECRET_KEY is found, and a new one is generated.
const SecretKeyWasNotSet = "SECRET_KEY was not set. A new one was generated for this runtime session." // #nosec G101

// Validation message constants.
const (
	ErrHTTPPortRequired               = "HTTP ports is required"
	ErrHTTPHostRequired               = "HTTP host is required"
	ErrHTTPContextPathEmpty           = "HTTP context path cannot be empty"
	ErrHTTPContextMustStart           = "HTTP context path must start with '/'"
	ErrHTTPContextMustNotEndWithSlash = "HTTP context path must not end with '/'"
	ErrHTTPReadTimeoutMin             = "HTTP read timeout must be at least %v"
	ErrHTTPWriteTimeoutMin            = "HTTP write timeout must be at least %v" // #nosec G101
	ErrHTTPReadHeaderTimeoutMin       = "HTTP read header timeout must be greater than 0"
	ErrHTTPIdleTimeoutMin             = "HTTP idle timeout must be greater than 0"
	ErrHTTPMaxHeaderBytesMin          = "HTTP max header bytes must be greater than 0"

	ErrGraphqlPortRequired         = "GraphQL ports is required"
	ErrGraphqlPathRequired         = "GraphQL path is required"
	ErrGraphqlPathMustStart        = "GraphQL path must start with '/'"
	ErrGraphqlReadHeaderTimeoutMin = "GraphQL read header timeout be greater than 0"
	ErrGraphqlReadTimeoutMin       = "GraphQL read timeout must be at least %v"
	ErrGraphqlWriteTimeoutMin      = "GraphQL write timeout must be at least %v"

	ErrCachePoolSizeMin = "CACHE_POOL_SIZE must be at least %d"
	ErrCacheAddrEmpty   = "cache address cannot be empty"

	ErrDBTypeEmpty          = "DB_TYPE cannot be empty"
	ErrDBHostEmpty          = "DB_HOST cannot be empty"
	ErrDBPortEmpty          = "DB_PORT cannot be empty"
	ErrDBNameRequired       = "database name is required"
	ErrDBUserRequired       = "database user is required"
	ErrDBPasswordRequired   = "database password is required"
	ErrDBTimeZoneEmpty      = "DB timezone (TZ) cannot be empty"
	ErrDBSSLModeInvalid     = "invalid DB_SSL_MODE value: %s"
	ErrDBMaxOpenConnsMin    = "DB_MAX_CONNECTIONS must be at least %d"
	ErrDBMaxIdleConnsNeg    = "DB_MAX_IDLE_CONNECTIONS cannot be negative"
	ErrDBConnMaxLifetimeNeg = "DB_CONN_MAX_LIFETIME_MINUTES cannot be negative"
	ErrDBRetryIntervalMin   = "DB_CONNECT_RETRY_INTERVAL must be at least %v"
	ErrDBMaxRetriesMin      = "DB_CONNECT_MAX_RETRIES must be at least %d"

	ErrOtelEndpointEmpty      = "OTel Exporter endpoint cannot be empty"
	ErrOtelCompressionInvalid = "OTel Exporter compression must be either 'none' or 'gzip', got: %s"

	ErrAppContextReqMin      = "context request timeout must be at least %v"
	ErrAppShutdownTimeoutMin = "shutdown timeout must be at least %s second"

	InfoSecretKeyGenerated = "JWT secret key successfully generated with length: %d" // #nosec G101
)
