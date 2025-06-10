// Package constants contains constants used throughout the application.
package constants

// MsgCacheConnected is a constant string used to log a message indicating that the cache has been successfully connected.
const MsgCacheConnected = "Cache connected"

// MsgPostgresConnected is a constant string used to indicate that the database connection has been successfully established.
const MsgPostgresConnected = "Database connection established"

// ErrConnectToCache indicates a failure to establish a connection with the Cache. The %v placeholder is used to include specific error details.
const ErrConnectToCache = "failed to connect to Cache: %v"

// ErrConnectToDatabase is an error message indicating a failure to establish a database connection. The placeholder %v is used for the specific error details.
const ErrConnectToDatabase = "Failed to connect to Database: %v"

// ErrCloseCacheConnection indicates an error occurred while attempting to close the Cache connection.
const ErrCloseCacheConnection = "failed to close Cache connection"

// FieldAddr is a constant used as a key for specifying address information in structured logging or configuration contexts.
const FieldAddr = "addr"
