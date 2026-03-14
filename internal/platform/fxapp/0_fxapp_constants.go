// Package fxapp provides constants for dependency injection wiring.
package fxapp

const (
	// Context names for cache initialization logging.
	contextNameAuth     = "auth"
	contextNameCategory = "category"
	contextNameTag      = "tag"
	contextNameRecord   = "record"
	contextNameUser     = "user"
	contextNameChat     = "chat"

	// Log messages for infrastructure lifecycle.
	logMsgConfigLoaded      = "configuration loaded"
	logMsgCacheInit         = "cache initialized"
	logMsgCacheCreateFailed = "failed to create cache"
	logMsgCacheAllClosed    = "all cache connections closed"
	logMsgDBConnectFailed   = "failed to connect to database"
	logMsgDBInitialized     = "database connection initialized"
	logMsgServerStarting    = "servers starting..."
	logMsgServerError       = "http server error"

	// Database type identifier.
	dbTypePostgresql = "postgresql"

	// Error message formatting.
	errMsgHTTPShutdown = "http shutdown: %w"
)
