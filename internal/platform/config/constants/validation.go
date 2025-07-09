package constants

import "time"

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
	MinShutdownTimeout = 1

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
