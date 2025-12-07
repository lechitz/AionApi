// Package usecase contains the business logic for the Chat context.
package usecase

// Tracing.
const (
	// TracerChatService is the name of the tracer for the chat service.
	TracerChatService = "chat.service"

	// SpanProcessMessage is the span name for processing a chat message.
	SpanProcessMessage = "ProcessMessage"
)

// Span attributes.
const (
	// AttrUserID is the attribute key for user ID.
	AttrUserID = "user_id"
	// AttrMessageLength is the attribute key for message length.
	AttrMessageLength = "message_length"
	// AttrTokensUsed is the attribute key for tokens used.
	AttrTokensUsed = "tokens_used"
	// AttrFunctionCallsCount is the attribute key for function calls count.
	AttrFunctionCallsCount = "function_calls_count"
)

// Span status messages.
const (
	// StatusMessageProcessedSuccessfully indicates a message was processed successfully.
	StatusMessageProcessedSuccessfully = "message processed successfully"
	// StatusFailedToCallAionChat indicates a failure to call the Aion-Chat service.
	StatusFailedToCallAionChat = "failed to call aion-chat service"
)

// Log messages.
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
