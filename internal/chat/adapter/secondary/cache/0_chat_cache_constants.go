// Package cache constants contains constants used for chat history cache operations.
package cache

import "time"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for chat cache operations.
// Format: aionapi.<domain>.<layer>.
const TracerName = "aionapi.chat.cache"

// -----------------------------------------------------------------------------
// Span Names
// -----------------------------------------------------------------------------

const (
	// SpanNameChatHistoryGet is the span name for retrieving chat history.
	SpanNameChatHistoryGet = "chat.cache.history.get"

	// SpanNameChatHistorySave is the span name for saving chat history.
	SpanNameChatHistorySave = "chat.cache.history.save"

	// SpanNameChatHistoryClear is the span name for clearing chat history.
	SpanNameChatHistoryClear = "chat.cache.history.clear"
)

// -----------------------------------------------------------------------------
// Attribute Keys
// -----------------------------------------------------------------------------

const (
	// AttributeCacheKey is the attribute key for cache keys.
	AttributeCacheKey = "cache_key"

	// AttributeCacheHit indicates cache hit/miss.
	AttributeCacheHit = "cache_hit"

	// AttributeLimit is the attribute key for history limit.
	AttributeLimit = "limit"

	// AttributeResultsCount is the attribute key for result count.
	AttributeResultsCount = "results_count"

	// AttributeHistorySize is the attribute key for history size.
	AttributeHistorySize = "history_size"
)

// =============================================================================
// BUSINESS LOGIC - Configuration and Messages.
// =============================================================================

const (
	// Cache key prefix for chat history.
	cacheKeyPrefix = "chat:history:"

	// Default TTL for chat history cache (30 minutes).
	defaultTTL = 30 * time.Minute

	// historyLimitDefault defines the max number of entries kept in cache.
	historyLimitDefault = 20
)

// Log messages.
const (
	// LogGetHistory is the log message for getting chat history from cache.
	LogGetHistory = "Getting chat history from cache"

	// LogHistoryRetrieved is the log message for successfully retrieved chat history.
	LogHistoryRetrieved = "Chat history retrieved from cache"

	// LogHistoryCacheMiss is the log message for cache miss.
	LogHistoryCacheMiss = "Chat history cache miss"

	// LogSetHistory is the log message for setting chat history in cache.
	LogSetHistory = "Setting chat history in cache"

	// LogHistorySet is the log message for successfully set chat history.
	LogHistorySet = "Chat history set in cache"

	// LogClearHistory is the log message for clearing chat history cache.
	LogClearHistory = "Clearing chat history cache"

	// LogHistoryCleared is the log message for successfully cleared chat history.
	LogHistoryCleared = "Chat history cache cleared"
)

// Error messages.
const (
	// ErrGetHistory is the error message for getting chat history from cache.
	ErrGetHistory = "failed to get chat history from cache"

	// ErrSetHistory is the error message for setting chat history in cache.
	ErrSetHistory = "failed to set chat history in cache"

	// ErrClearHistory is the error message for clearing chat history cache.
	ErrClearHistory = "failed to clear chat history cache"

	// ErrUnmarshalHistory is the error message for unmarshaling chat history.
	ErrUnmarshalHistory = "failed to unmarshal chat history"

	// ErrMarshalHistory is the error message for marshaling chat history.
	ErrMarshalHistory = "failed to marshal chat history"
)
