// Package handler implements HTTP handlers for chat endpoints.
package handler

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerChatHandler is the tracer name for chat handler.
// Format: aionapi.<domain>.<layer> .
const TracerChatHandler = "aionapi.chat.handler"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanChatHandler is the span name for handling chat requests.
	SpanChatHandler = "chat.handler.handle"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// EventDecodeRequest indicates the decoding of the chat request.
	EventDecodeRequest = "chat.handler.decode_request"
	// EventValidateRequest indicates the validation of the chat request.
	EventValidateRequest = "chat.handler.validate_request"
	// EventCallService indicates the call to the chat service.
	EventCallService = "chat.handler.call_service"
	// EventChatSuccess indicates successful chat processing.
	EventChatSuccess = "chat.handler.success"
	// EventChatError indicates an error occurred during chat processing.
	EventChatError = "chat.handler.error"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusChatSuccess indicates successful chat processing.
	StatusChatSuccess = "chat processed successfully"
)

// =============================================================================
// BUSINESS LOGIC - Error and Success Messages
// =============================================================================

// Error messages.
const (
	// ErrChat indicates a failure in processing the chat message.
	ErrChat = "failed to process chat message"
	// ErrRequiredMessage indicates that the message field is required.
	ErrRequiredMessage = "message is required"
	// ErrMessageTooShort indicates that the message is too short.
	ErrMessageTooShort = "message must be at least 1 character"
	// ErrMessageTooLong indicates that the message exceeds the maximum length.
	ErrMessageTooLong = "message must not exceed 2000 characters"
	// ErrUserIDNotFound indicates that the user ID was not found in the context.
	ErrUserIDNotFound = "user ID not found in context"
)

// Log messages.
const (
	MsgChatSuccess      = "Chat processed successfully"
	MsgChatRequestStart = "Processing chat request"
)

// Validation constraints.
const (
	// MinMessageLength is the minimum length for a chat message.
	MinMessageLength = 1
	// MaxMessageLength is the maximum length for a chat message.
	MaxMessageLength = 2000
)
