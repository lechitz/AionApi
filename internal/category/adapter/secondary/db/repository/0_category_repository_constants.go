// Package repository defines DB-facing repository logic for Category.
// This file centralizes all magic strings/constants used by the repository.
package repository

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the tracer name used by the category repository.
// Format: aion-api.<domain>.<layer>.
const TracerName = "aion-api.category.repository"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreateRepo is the span name for creating a category.
	SpanCreateRepo = "category.repository.create"

	// SpanGetByIDRepo is the span name for fetching a category by ID.
	SpanGetByIDRepo = "category.repository.get_by_id"

	// SpanGetByNameRepo is the span name for fetching a category by name.
	SpanGetByNameRepo = "category.repository.get_by_name"

	// SpanListAllRepo is the span name for listing all categories.
	SpanListAllRepo = "category.repository.list_all"

	// SpanSoftDeleteRepo is the span name for soft-deleting a category.
	SpanSoftDeleteRepo = "category.repository.soft_delete"

	// SpanUpdateRepo is the span name for updating a category.
	SpanUpdateRepo = "category.repository.update"
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

	// OpGetByName is the attribute value for the "get by name" operation.
	OpGetByName = "get_by_name"

	// OpListAll is the attribute value for the "list all" operation.
	OpListAll = "list_all"

	// OpSoftDelete is the attribute value for the "soft delete" operation.
	OpSoftDelete = "soft_delete"

	// OpUpdate is the attribute value for the "update" operation.
	OpUpdate = "update"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusCategoryCreated indicates the category was created successfully.
	StatusCategoryCreated = "category created successfully"

	// StatusRetrievedByID indicates a category was retrieved by ID successfully.
	StatusRetrievedByID = "category retrieved by id successfully"

	// StatusRetrievedByName indicates a category was retrieved by Name successfully.
	StatusRetrievedByName = "category retrieved by name successfully"

	// StatusFetchedAll indicates all categories were retrieved successfully.
	StatusFetchedAll = "all categories retrieved successfully"

	// StatusSoftDeleted indicates a category was soft-deleted successfully.
	StatusSoftDeleted = "category soft deleted successfully"

	// StatusUpdated indicates a category was updated successfully.
	StatusUpdated = "category updated successfully"
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrCreateCategoryMsg is the error message used when creation fails.
	ErrCreateCategoryMsg = "error creating category"

	// ErrGetCategoryMsg is the error message used when a generic get fails.
	ErrGetCategoryMsg = "error getting category"

	// ErrCategoryNotFoundMsg is the error message used when a category is not found.
	ErrCategoryNotFoundMsg = "category not found"
)
