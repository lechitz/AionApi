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
	ErrAuditService       = "failed to list audit events"
	ErrMissingUserID      = "user ID not found in context"
	ErrInvalidUserID      = "invalid user ID"
	ErrForbiddenCrossUser = "forbidden cross-user audit query"
)

const (
	MsgAuditEventsListed = "Audit events listed successfully"
	// MsgCrossUserAuditQuery is emitted when an admin queries another user's audit events.
	MsgCrossUserAuditQuery = "cross-user audit query"
	// MsgCrossUserAuditQueryFailed is emitted when cross-user query fails.
	MsgCrossUserAuditQueryFailed = "cross-user audit query failed"
)

const (
	QueryTraceID = "trace_id"
	QueryDraftID = "draft_id"
	QueryStatus  = "status"
	QueryUserID  = "user_id"
	QueryFromUTC = "from_utc"
	QueryToUTC   = "to_utc"
	QueryLimit   = "limit"
	QueryOffset  = "offset"
)

const (
	DefaultLimit = 100
	MaxLimit     = 500
)
