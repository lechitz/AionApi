// Package cache constants contains constants used for token operations.
package cache

import "time"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for token cache operations.
// Format: aionapi.<domain>.<layer>.
const TracerName = "aionapi.auth.cache"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanNameTokenSave is the span name for saving a token.
	SpanNameTokenSave = "auth.cache.token_save" //nolint:gosec // span name, not a credential

	// SpanNameTokenGet is the span name for retrieving a token.
	SpanNameTokenGet = "auth.cache.token_get" //nolint:gosec // span name, not a credential

	// SpanNameTokenDelete is the span name for deleting a token.
	SpanNameTokenDelete = "auth.cache.token_delete" //nolint:gosec // span name, not a credential

	// SpanNameTokenSaveWithKey is the span name for saving a token with custom key.
	SpanNameTokenSaveWithKey = "auth.cache.token_save_with_key" //nolint:gosec // span name, not a credential

	// SpanNameTokenGetByKey is the span name for retrieving a token by custom key.
	SpanNameTokenGetByKey = "auth.cache.token_get_by_key" //nolint:gosec // span name, not a credential
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// OperationSave is the name of the save operation.
	OperationSave = "save"

	// OperationGet is the name of the get operation.
	OperationGet = "get"

	// OperationDelete is the name of the delete operation.
	OperationDelete = "delete"
)

// Attribute keys for tracing and logging.
const (
	// AttributeCacheKey is the attribute key for cache keys.
	AttributeCacheKey = "cache_key"

	// AttributeTTL is the attribute key for time-to-live values.
	AttributeTTL = "ttl"
)

// =============================================================================
// BUSINESS LOGIC - Configuration and Messages
// =============================================================================

// TokenExpirationDefault defines the default expiration for tokens.
const TokenExpirationDefault = 24 * time.Hour

// TokenUserKeyFormat defines the format for token keys associated with a specific user and token type.
//
//nolint:gosec // Cache key pattern for Redis; contains the word "token" but is not a credential/secret.
const TokenUserKeyFormat = "token:user:%d:%s"

// Success messages.
const (
	// TokenRetrievedSuccessfully is the message used when a token is retrieved successfully.
	TokenRetrievedSuccessfully = "token retrieved successfully"

	// TokenDeletedSuccessfully is the message used when a token is deleted successfully.
	TokenDeletedSuccessfully = "token deleted successfully"

	// TokenSavedSuccessfully is the message used when a token is saved successfully.
	TokenSavedSuccessfully = "token saved successfully" // #nosec G101 - false positive, this is a log message not a credential

	// TokenRetrievedByCustomKey is the message used when a token is retrieved by custom key.
	TokenRetrievedByCustomKey = "token retrieved by custom key"

	// TokenSavedWithCustomKey is the message used when a token is saved with custom key.
	TokenSavedWithCustomKey = "token saved with custom key"
)

// Error messages.
const (
	// ErrorToSaveTokenToRedis indicates a failure to save a token in Redis.
	ErrorToSaveTokenToRedis = "error to save token to Redis"

	// ErrorToGetTokenFromRedis indicates a failure to retrieve a token from Redis.
	ErrorToGetTokenFromRedis = "error to get token from Redis"

	// ErrorToDeleteTokenFromRedis indicates a failure to delete a token from Redis.
	ErrorToDeleteTokenFromRedis = "error to delete token from Redis"

	// ErrorTokenNotFoundInGracePeriod indicates a token was not found in grace period cache.
	ErrorTokenNotFoundInGracePeriod = "token not found in grace period"
)
