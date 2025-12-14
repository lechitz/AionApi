// Package repository defines DB-facing repository logic for Tag.
// This file centralizes all magic strings/constants used by the repository.
package repository

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the tracer name used by the Tag repository.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.tag.repository"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreateRepo is the span name for creating a Tag.
	SpanCreateRepo = "tag.repository.create"

	// SpanGetByNameRepo is the span name for fetching a Tag by name.
	SpanGetByNameRepo = "tag.repository.get_by_name"

	// SpanGetAllRepo is the span name for fetching all Tags for a user.
	SpanGetAllRepo = "tag.repository.get_all"

	// SpanGetByCategoryRepo is the span name for fetching Tags by category.
	SpanGetByCategoryRepo = "tag.repository.get_by_category"

	// SpanSoftDeleteRepo is the span name for soft-deleting a Tag.
	SpanSoftDeleteRepo = "tag.repository.soft_delete"
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// OpCreate is the attribute value for the "create" operation.
	OpCreate = "create"

	// OpGetByName is the attribute value for the "get by name" operation.
	OpGetByName = "get_by_name"

	// OpGetByID is the attribute value for the "get by id" operation.
	OpGetByID = "get_by_id"

	// OpGetAll is the attribute value for the "get all" operation.
	OpGetAll = "get_all"

	// OpGetByCategory is the attribute value for the "get by category" operation.
	OpGetByCategory = "get_by_category"

	// OpSoftDelete is the attribute value for the "soft delete" operation.
	OpSoftDelete = "soft_delete"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusTagCreated indicates the Tag was created successfully.
	StatusTagCreated = "tag created successfully"

	// StatusRetrievedByID indicates a Tag was retrieved by ID successfully.
	StatusRetrievedByID = "tag retrieved by id successfully"

	// StatusRetrievedAll indicates all Tags were retrieved successfully.
	StatusRetrievedAll = "tags retrieved successfully"

	// StatusRetrievedByCategory indicates Tags were retrieved by category successfully.
	StatusRetrievedByCategory = "tags retrieved by category successfully"

	// StatusRetrievedByName indicates a Tag was retrieved by Name successfully.
	StatusRetrievedByName = "tag retrieved by name successfully"

	// StatusFetchedAll indicates all categories were retrieved successfully.
	StatusFetchedAll = "all categories retrieved successfully"

	// StatusSoftDeleted indicates a Tag was soft-deleted successfully.
	StatusSoftDeleted = "tag soft deleted successfully"

	// StatusUpdated indicates a Tag was updated successfully.
	StatusUpdated = "tag updated successfully"
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrCreateTagMsg is the error message used when creation fails.
	ErrCreateTagMsg = "error creating tag"

	// ErrGetTagByIDMsg is the error message used when a generic get fails.
	ErrGetTagByIDMsg = "error getting tag by ID"

	// ErrGetAllTagsMsg is the error message used when getting all Tags fails.
	ErrGetAllTagsMsg = "error getting all tags for user"

	// ErrGetTagsByCategoryMsg is the error message used when getting Tags by category fails.
	ErrGetTagsByCategoryMsg = "error getting tags by category"

	// ErrGetTagByNameMsg is the error message used when a generic get fails.
	ErrGetTagByNameMsg = "error getting tag by name"

	// ErrTagNotFoundMsg is the error message used when a Tag is not found.
	ErrTagNotFoundMsg = "tag not found"

	// ErrNoTagsFoundMsg is the error message used when no Tags are found for a user.
	ErrNoTagsFoundMsg = "no tags found for user"
)
