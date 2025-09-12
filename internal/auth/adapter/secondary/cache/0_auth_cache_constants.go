// Package cache constants contains constants used for token operations.
package cache

import "time"

// ErrorToSaveTokenToRedis indicates a failure to save a token in Redis.
const ErrorToSaveTokenToRedis = "error to save token to Redis"

// ErrorToGetTokenFromRedis indicates a failure to retrieve a token from Redis.
const ErrorToGetTokenFromRedis = "error to get token from Redis"

// ErrorToDeleteTokenFromRedis indicates a failure to delete a token from Redis.
const ErrorToDeleteTokenFromRedis = "error to delete token from Redis"

// TokenExpirationDefault defines the default expiration for tokens.
const TokenExpirationDefault = 24 * time.Hour

// TokenUserKeyFormat defines the format for token keys associated with a specific user.
//
//nolint:gosec // Cache key pattern for Redis; contains the word "token" but is not a credential/secret.
const TokenUserKeyFormat = "token:user:%d"

const (
	// TokenRetrievedSuccessfully is the message used when a token is retrieved successfully.
	TokenRetrievedSuccessfully = "token retrieved successfully"

	// TokenDeletedSuccessfully is the message used when a token is deleted successfully.
	TokenDeletedSuccessfully = "token deleted successfully"

	// TokenSavedSuccessfully is the message used when a token is saved successfully.
	//nolint:gosec // Static log/message string; contains the word "token" but is not a credential/secret.
	TokenSavedSuccessfully = "token saved successfully"
)

const (
	// SpanTracerTokenStore is the name of the tracer for token operations.
	SpanTracerTokenStore = "TokenStore"

	// SpanNameTokenSave is the name of the span for saving a token.
	SpanNameTokenSave = "token.save"

	// SpanNameTokenGet is the name of the span for retrieving a token.
	SpanNameTokenGet = "token.get"

	// SpanNameTokenDelete is the name of the span for deleting a token.
	SpanNameTokenDelete = "token.delete"
)
