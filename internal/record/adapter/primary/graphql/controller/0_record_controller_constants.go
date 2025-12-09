// Package controller contains GraphQL-facing controllers for the Record context.
package controller

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for record GraphQL controllers.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.record.controller"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreate is the span name for creating a record.
	SpanCreate = "record.controller.create"

	// SpanGetByName is the span name for retrieving a record by name.
	SpanGetByName = "record.controller.get_by_name"

	// SpanGetByCategory is the span name for retrieving records by category.
	SpanGetByCategory = "record.controller.get_by_category"

	// SpanGetByTag is the span name for retrieving a record by tag.
	SpanGetByTag = "record.controller.get_by_tag"

	// SpanListAll is the span name for listing all records for a user.
	SpanListAll = "record.controller.list_all"

	// SpanListByTag is the span name for listing records by tag.
	SpanListByTag = "record.controller.list_by_tag"

	// SpanListByDay is the span name for listing records by day.
	SpanListByDay = "record.controller.list_by_day"

	// SpanListAllUntil is the span name for listing records until a timestamp.
	SpanListAllUntil = "record.controller.list_all_until"

	// SpanListAllBetween is the span name for listing records between dates.
	SpanListAllBetween = "record.controller.list_all_between"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// MsgCreated is the log message for when a record is created.
	MsgCreated = "record created"

	// StatusCreated is the human-readable status used in responses/logs when a record is created.
	StatusCreated = "created"

	// StatusFetched is the human-readable status used in responses/logs when records are retrieved.
	StatusFetched = "fetched"
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// MsgCreateError is the log message for when a create operation fails.
	MsgCreateError = "error creating record"

	// ErrUserIDNotFound is the error message when the user ID is missing or invalid.
	ErrUserIDNotFound = "user id not found"

	// ErrRecordNotFound is the error message when the record ID is missing or invalid.
	ErrRecordNotFound = "record id not found"

	// ErrInvalidRecordID is the error message when the record ID cannot be parsed or is invalid.
	ErrInvalidRecordID = "invalid record id"

	// ErrInvalidCategoryID is the error message when the category ID cannot be parsed or is invalid.
	ErrInvalidCategoryID = "invalid record id"

	// ErrInvalidTagID is the error message when the tag ID cannot be parsed or is invalid.
	ErrInvalidTagID = "invalid tag id"

	// ErrTagNotFound is the error message when a record cannot be found.
	ErrTagNotFound = "tag not found"

	// ErrFailedToListRecords is the error message when listing records fails.
	ErrFailedToListRecords = "failed to list records"
)
