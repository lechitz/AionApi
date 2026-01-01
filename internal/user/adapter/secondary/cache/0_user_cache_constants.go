// Package cache provides user cache operations.
package cache

import "time"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for user cache operations.
const TracerName = "aionapi.user.cache"

// Span Names.
const (
	SpanNameUserSave   = "user.cache.save"
	SpanNameUserGet    = "user.cache.get"
	SpanNameUserDelete = "user.cache.delete"
)

// Operations.
const (
	OperationSave   = "save"
	OperationGet    = "get"
	OperationDelete = "delete"
)

// Attribute keys for tracing and logging.
const (
	AttributeCacheKey = "cache_key"
	AttributeTTL      = "ttl"
)

// =============================================================================
// BUSINESS LOGIC - Configuration and Messages
// =============================================================================

// UserExpirationDefault defines the default expiration for user profile cache (1 hour).
// User data is relatively stable, so we can cache longer than categories/tags.
const UserExpirationDefault = 60 * time.Minute

// Key formats.
const (
	// UserIDKeyFormat : user:id:{userID}.
	UserIDKeyFormat = "user:id:%d"

	// UserUsernameKeyFormat : user:username:{username}.
	UserUsernameKeyFormat = "user:username:%s"

	// UserEmailKeyFormat : user:email:{email}.
	UserEmailKeyFormat = "user:email:%s"
)

// Success messages.
const (
	UserRetrievedSuccessfully = "user retrieved successfully from cache"
	UserDeletedSuccessfully   = "user deleted successfully from cache"
	UserSavedSuccessfully     = "user saved successfully to cache" // #nosec G101 - false positive, this is a log message
)

// Error messages.
const (
	ErrorToSaveUserToCache     = "error to save user to cache"
	ErrorToGetUserFromCache    = "error to get user from cache"
	ErrorToDeleteUserFromCache = "error to delete user from cache"
	ErrorToSerializeUser       = "error to serialize user"
	ErrorToDeserializeUser     = "error to deserialize user"
)
