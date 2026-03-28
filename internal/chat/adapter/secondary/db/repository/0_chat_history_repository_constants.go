// Package repository defines DB-facing repository logic for ChatHistory.
// This file centralizes all magic strings/constants used by the repository.
package repository

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the tracer name used by the ChatHistory repository.
// Format: aion-api.<domain>.<layer> .
const TracerName = "aion-api.chat.repository"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanSaveRepo is the span name for saving a ChatHistory entry.
	SpanSaveRepo = "chat.repository.save"

	// SpanGetLatestRepo is the span name for fetching latest chat history.
	SpanGetLatestRepo = "chat.repository.get_latest"

	// SpanGetByUserIDRepo is the span name for fetching chat history by user ID.
	SpanGetByUserIDRepo = "chat.repository.get_by_user_id"
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// OpSave is the attribute value for the "save" operation.
	OpSave = "save"

	// OpGetLatest is the attribute value for the "get latest" operation.
	OpGetLatest = "get_latest"

	// OpGetByUserID is the attribute value for the "get by user id" operation.
	OpGetByUserID = "get_by_user_id"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusChatSaved indicates the chat history was saved successfully.
	StatusChatSaved = "chat history saved successfully"

	// StatusRetrievedLatest indicates latest chat history was retrieved successfully.
	StatusRetrievedLatest = "latest chat history retrieved successfully"

	// StatusRetrievedByUserID indicates chat history was retrieved by user ID successfully.
	StatusRetrievedByUserID = "chat history retrieved by user ID successfully"
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrSaveChatMsg is the error message used when saving chat history fails.
	ErrSaveChatMsg = "error saving chat history"

	// ErrGetLatestChatMsg is the error message used when getting latest chat history fails.
	ErrGetLatestChatMsg = "error getting latest chat history"

	// ErrGetChatByUserIDMsg is the error message used when getting chat history by user ID fails.
	ErrGetChatByUserIDMsg = "error getting chat history by user ID"

	// ErrNoChatHistoryFoundMsg is the error message used when no chat history is found.
	ErrNoChatHistoryFoundMsg = "no chat history found for user"
)
