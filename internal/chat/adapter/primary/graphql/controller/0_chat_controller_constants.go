// Package controller provides constants for the Chat GraphQL controller.
package controller

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for Chat GraphQL controllers.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.chat.controller"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanGetChatHistory is the span name for getting chat history.
	SpanGetChatHistory = "chat.controller.get_history"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusFetched indicates that chat history has been successfully fetched.
	StatusFetched = "chat_history_fetched"
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
)

// -----------------------------------------------------------------------------
// Sentinel Errors
// -----------------------------------------------------------------------------

// (No specific errors for now, using service errors)
