// Package usecase contains business logic for the audit context.
package usecase

const (
	// TracerName is the tracer name for audit use cases.
	TracerName = "aionapi.audit.usecase"
)

const (
	// SpanWriteEvent is the span name for writing one audit event.
	SpanWriteEvent = "audit.event.write"
	// SpanListEvents is the span name for listing audit events.
	SpanListEvents = "audit.event.list"
)

const (
	// StatusAuditEventWritten indicates event persisted successfully.
	StatusAuditEventWritten = "audit event written successfully"
	// StatusAuditEventsListed indicates events queried successfully.
	StatusAuditEventsListed = "audit events listed successfully"
)

const (
	// LogWritingAuditEvent indicates audit write start.
	LogWritingAuditEvent = "Writing audit event"
	// LogAuditEventWritten indicates audit write completion.
	LogAuditEventWritten = "Audit event written"
	// LogFailedWriteAuditEvent indicates audit write failure.
	LogFailedWriteAuditEvent = "Failed to write audit event"
	// LogListingAuditEvents indicates audit listing start.
	LogListingAuditEvents = "Listing audit events"
	// LogFailedListAuditEvents indicates audit list failure.
	LogFailedListAuditEvents = "Failed to list audit events"
)

const (
	// LogKeyError is the generic error key.
	LogKeyError = "error"
	// LogKeyUserID is the key for user identifier.
	LogKeyUserID = "user_id"
	// LogKeyTraceID is the key for trace identifier.
	LogKeyTraceID = "trace_id"
	// LogKeyDraftID is the key for draft identifier.
	LogKeyDraftID = "draft_id"
)
