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
	//nolint:gosec // Static log/message string; contains the word "token" but is not a credential/secret.
	TokenSavedSuccessfully = "token saved successfully"
)

// Error messages.
const (
	// ErrorToSaveTokenToRedis indicates a failure to save a token in Redis.
	ErrorToSaveTokenToRedis = "error to save token to Redis"

	// ErrorToGetTokenFromRedis indicates a failure to retrieve a token from Redis.
	ErrorToGetTokenFromRedis = "error to get token from Redis"

	// ErrorToDeleteTokenFromRedis indicates a failure to delete a token from Redis.
	ErrorToDeleteTokenFromRedis = "error to delete token from Redis"
)
