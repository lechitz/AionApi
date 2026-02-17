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

	// SpanSaveChatHistory is the span name for saving chat history.
	SpanSaveChatHistory = "chat.history.save"

	// SpanGetChatHistory is the span name for getting chat history.
	SpanGetChatHistory = "chat.history.get"

	// SpanGetChatContext is the span name for getting chat context.
	SpanGetChatContext = "chat.context.get"
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
	// StatusChatHistorySaved indicates chat history was saved successfully.
	StatusChatHistorySaved = "chat history saved successfully"
	// StatusChatHistoryRetrieved indicates chat history was retrieved successfully.
	StatusChatHistoryRetrieved = "chat history retrieved successfully"
	// StatusChatContextRetrieved indicates chat context was retrieved successfully.
	StatusChatContextRetrieved = "chat context retrieved successfully"
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
	// LogSavingChatHistory is the log message for saving chat history.
	LogSavingChatHistory = "Saving chat history"
	// LogChatHistorySaved is the log message for successfully saved chat history.
	LogChatHistorySaved = "Chat history saved successfully"
	// LogFailedToSaveChatHistory is the log message for failed to save chat history.
	LogFailedToSaveChatHistory = "Failed to save chat history"
	// LogGettingChatHistory is the log message for getting chat history.
	LogGettingChatHistory = "Getting chat history"
	// LogChatHistoryRetrieved is the log message for successfully retrieved chat history.
	LogChatHistoryRetrieved = "Chat history retrieved successfully"
	// LogFailedToGetChatHistory is the log message for failed to get chat history.
	LogFailedToGetChatHistory = "Failed to get chat history"
	// LogGettingChatContext is the log message for getting chat context.
	LogGettingChatContext = "Getting chat context"
	// LogChatContextRetrieved is the log message for successfully retrieved chat context.
	LogChatContextRetrieved = "Chat context retrieved successfully"
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
	// ContextKeyUserTimezone is the canonical context key for timezone sent to aion-chat.
	ContextKeyUserTimezone = "user_timezone"
	// DefaultTimezone is the default timezone value.
	DefaultTimezone = "America/Sao_Paulo"
)
