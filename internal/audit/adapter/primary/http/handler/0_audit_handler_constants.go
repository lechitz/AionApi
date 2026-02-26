// Package handler implements HTTP handlers for audit endpoints.
package handler

const (
	// TracerAuditHandler is the tracer name for audit HTTP handlers.
	TracerAuditHandler = "aionapi.audit.handler"
)

const (
	// SpanListEventsHandler is the span name for listing audit events.
	SpanListEventsHandler = "audit.handler.list_events"
)

const (
	errAuditService       = "failed to list audit events"
	errMissingUserID      = "user ID not found in context"
	errInvalidUserID      = "invalid user ID"
	errForbiddenCrossUser = "forbidden cross-user audit query"
)

const (
	msgAuditEventsListed = "Audit events listed successfully"
	// MsgCrossUserAuditQuery is emitted when an admin queries another user's audit events.
	MsgCrossUserAuditQuery = "cross-user audit query"
	// MsgCrossUserAuditQueryFailed is emitted when cross-user query fails.
	MsgCrossUserAuditQueryFailed = "cross-user audit query failed"
)

const (
	queryTraceID = "trace_id"
	queryDraftID = "draft_id"
	queryStatus  = "status"
	queryUserID  = "user_id"
	queryFromUTC = "from_utc"
	queryToUTC   = "to_utc"
	queryLimit   = "limit"
	queryOffset  = "offset"
)

const (
	defaultLimit = 100
	maxLimit     = 500
)
