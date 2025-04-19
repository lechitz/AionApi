package constants

const (
	// Connection Messages
	MsgCacheConnected    = "Cache connected"
	MsgPostgresConnected = "Database connection established"

	// Error Messages
	ErrConnectToCache       = "failed to connect to Cache: %v"
	ErrConnectToDatabase    = "Failed to connect to Database: %v"
	ErrCloseCacheConnection = "failed to close Cache connection"

	// Log Fields
	FieldAddr = "addr"
)
