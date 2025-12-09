// Package repository defines DB-facing repository logic for Record.
package repository

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the tracer name used by the Record repository.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.record.repository"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreateRepo is the span name for creating a record.
	SpanCreateRepo = "record.repository.create"

	// SpanGetByIDRepo is the span name for fetching a record by ID.
	SpanGetByIDRepo = "record.repository.get_by_id"

	// SpanGetByTagRepo is the span name for fetching records by tag.
	SpanGetByTagRepo = "record.repository.get_by_tag"

	// SpanListAllRepo is the span name for listing all records.
	SpanListAllRepo = "record.repository.list_all"

	// SpanListByTagRepo is the span name for listing records by tag.
	SpanListByTagRepo = "record.repository.list_by_tag"

	// SpanListByDayRepo is the span name for listing records by day.
	SpanListByDayRepo = "record.repository.list_by_day"

	// SpanUpdateRepo is the span name for updating a record.
	SpanUpdateRepo = "record.repository.update"

	// SpanSoftDeleteRepo is the span name for soft-deleting a record.
	SpanSoftDeleteRepo = "record.repository.soft_delete"
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// OpCreate is the attribute value for the "create" operation.
	OpCreate = "create"

	// OpGetByID is the attribute value for the "get by id" operation.
	OpGetByID = "get_by_id"

	// OpGetByTag is the attribute value for the "get by tag" operation.
	OpGetByTag = "get_by_tag"

	// OpListAll is the attribute value for the "list all" operation.
	OpListAll = "list_all"

	// OpListByTag is the attribute value for the "list by tag" operation.
	OpListByTag = "list_by_tag"

	// OpListByDay is the attribute value for the "list by day" operation.
	OpListByDay = "list_by_day"

	// OpUpdate is the attribute value for the "update" operation.
	OpUpdate = "update"

	// OpSoftDelete is the attribute value for the "soft delete" operation.
	OpSoftDelete = "soft_delete"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusRecordCreated indicates the record was created successfully.
	StatusRecordCreated = "record created successfully"

	// StatusRetrievedByID indicates a record was retrieved by ID successfully.
	StatusRetrievedByID = "record retrieved by id successfully"

	// StatusRetrievedByTag indicates records were retrieved by tag successfully.
	StatusRetrievedByTag = "records retrieved by tag successfully"

	// StatusFetchedAll indicates all records were retrieved successfully.
	StatusFetchedAll = "all records retrieved successfully"

	// StatusSoftDeleted indicates a record was soft-deleted successfully.
	StatusSoftDeleted = "record soft deleted successfully"

	// StatusUpdated indicates a record was updated successfully.
	StatusUpdated = "record updated successfully"
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrCreateRecordMsg is the error message used when creation fails.
	ErrCreateRecordMsg = "error creating record"

	// ErrGetRecordMsg is the error message used when a generic get fails.
	ErrGetRecordMsg = "error getting record"

	// ErrRecordNotFoundMsg is the error message used when a record is not found.
	ErrRecordNotFoundMsg = "record not found"

	// ErrListRecordsMsg is the error message used when listing records fails.
	ErrListRecordsMsg = "error listing records"

	// ErrUpdateRecordMsg is the error message used when updating a record fails.
	ErrUpdateRecordMsg = "error updating record"

	// ErrDeleteRecordMsg is the error message used when deleting a record fails.
	ErrDeleteRecordMsg = "error deleting record"
)
