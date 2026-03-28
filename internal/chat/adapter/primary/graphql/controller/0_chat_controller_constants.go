// Package controller provides constants for the Chat GraphQL controller.
package controller

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for Chat GraphQL controllers.
// Format: aion-api.<domain>.<layer> .
const TracerName = "aion-api.chat.controller"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanGetChatHistory is the span name for getting chat history.
	SpanGetChatHistory = "chat.controller.get_history"

	// SpanGetChatContext is the span name for getting chat context.
	SpanGetChatContext = "chat.controller.get_context"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusFetched indicates that chat history has been successfully fetched.
	StatusFetched = "chat_history_fetched"

	// StatusContextFetched indicates that chat context has been successfully fetched.
	StatusContextFetched = "chat_context_fetched"
)

// =============================================================================
// BUSINESS LOGIC - Error and Log Messages
// =============================================================================

// -----------------------------------------------------------------------------
// Log Messages
// -----------------------------------------------------------------------------

const (
	// MsgFetched is the log message for when chat history is fetched.
	MsgFetched = "chat history fetched"

	// MsgFetchError is the log message for when a fetch operation fails.
	MsgFetchError = "error fetching chat history"

	// MsgContextFetched is the log message for when chat context is fetched.
	MsgContextFetched = "chat context fetched"

	// MsgContextFetchError is the log message for when chat context fetch fails.
	MsgContextFetchError = "error fetching chat context"
)

// -----------------------------------------------------------------------------
// Attribute Keys
// Used in span.SetAttributes() and log metadata
// -----------------------------------------------------------------------------

const (
	// AttrLimit is the attribute key for pagination limit.
	AttrLimit = "limit"

	// AttrOffset is the attribute key for pagination offset.
	AttrOffset = "offset"

	// AttrCount is the attribute key for result count.
	AttrCount = "count"

	// AttrRecentChats is the attribute key for recent chats count in logs.
	AttrRecentChats = "recent_chats"

	// AttrRecentChatsCount is the attribute key for recent chats count in spans.
	AttrRecentChatsCount = "recent_chats_count"
)

// -----------------------------------------------------------------------------
// Sentinel Errors
// -----------------------------------------------------------------------------

// (No specific errors for now, using service errors)
