// Package usecase contains the business logic for the Chat context.
package usecase

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for the chat service.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.chat.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanProcessMessage is the span name for processing a chat message.
	SpanProcessMessage = "chat.message.process"
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// AttrUserID is the attribute key for user ID.
	AttrUserID = "aion.user_id"
	// AttrMessageLength is the attribute key for message length.
	AttrMessageLength = "aion.chat.message_length"
	// AttrTokensUsed is the attribute key for tokens used.
	AttrTokensUsed = "aion.chat.tokens_used"
	// AttrFunctionCallsCount is the attribute key for function calls count.
	AttrFunctionCallsCount = "aion.chat.function_calls_count"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusMessageProcessedSuccessfully indicates a message was processed successfully.
	StatusMessageProcessedSuccessfully = "message processed successfully"
	// StatusFailedToCallAionChat indicates a failure to call the Aion-Chat service.
	StatusFailedToCallAionChat = "failed to call aion-chat service"
)

// =============================================================================
// BUSINESS LOGIC - Log Messages
// =============================================================================

const (
	// LogProcessingChatMessage is the log message for processing a chat message.
	LogProcessingChatMessage = "Processing chat message"
	// LogChatMessageProcessedSuccessfully is the log message for a successfully processed message.
	LogChatMessageProcessedSuccessfully = "Chat message processed successfully"
	// LogFailedToCallAionChat is the log message for a failed call to Aion-Chat.
	LogFailedToCallAionChat = "Failed to call Aion-Chat service"
)

// Log keys.
const (
	// LogKeyUserID is the log key for user ID.
	LogKeyUserID = "user_id"
	// LogKeyMessageLength is the log key for message length.
	LogKeyMessageLength = "message_length"
	// LogKeyTokensUsed is the log key for tokens used.
	LogKeyTokensUsed = "tokens_used"
	// LogKeyResponseLength is the log key for response length.
	LogKeyResponseLength = "response_length"
	// LogKeyError is the log key for error.
	LogKeyError = "error"
)

// Context keys.
const (
	// ContextKeyTimezone is the context key for timezone.
	ContextKeyTimezone = "timezone"
	// DefaultTimezone is the default timezone value.
	DefaultTimezone = "America/Sao_Paulo"
)
