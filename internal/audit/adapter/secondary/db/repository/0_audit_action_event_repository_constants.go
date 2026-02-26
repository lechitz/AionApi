// Package repository implements DB repositories for audit persistence.
package repository

const (
	// AuditTracerName is the tracer name used by the audit repository.
	AuditTracerName = "aionapi.audit.repository"
)

const (
	// SpanAuditSaveRepo is the span name for persisting one audit event.
	SpanAuditSaveRepo = "audit.repository.save"
	// SpanAuditListRepo is the span name for listing audit events.
	SpanAuditListRepo = "audit.repository.list"
)

const (
	// OpAuditSave is the operation value for audit save.
	OpAuditSave = "audit_save"
	// OpAuditList is the operation value for audit list.
	OpAuditList = "audit_list"
)

const (
	// StatusAuditSaved indicates successful audit persistence.
	StatusAuditSaved = "audit action event saved successfully"
	// StatusAuditListed indicates successful audit retrieval.
	StatusAuditListed = "audit action events listed successfully"
)

const (
	// ErrSaveAuditActionEventMsg is used when saving audit event fails.
	ErrSaveAuditActionEventMsg = "error saving audit action event"
	// ErrListAuditActionEventMsg is used when listing audit events fails.
	ErrListAuditActionEventMsg = "error listing audit action events"
)
