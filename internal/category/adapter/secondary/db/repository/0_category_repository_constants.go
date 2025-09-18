// Package repository defines DB-facing repository logic for Category.
// This file centralizes all magic strings/constants used by the repository.
package repository

// Tracer & span names (OpenTelemetry)

// TracerName is the tracer name used by the category repository.
const TracerName = "CategoryRepository"

// SpanCreateRepo is the span name for creating a category.
const SpanCreateRepo = "Create"

// SpanGetByIDRepo is the span name for fetching a category by ID.
const SpanGetByIDRepo = "GetByID"

// SpanGetByNameRepo is the span name for fetching a category by name.
const SpanGetByNameRepo = "GetByName"

// SpanListAllRepo is the span name for listing all categories.
const SpanListAllRepo = "ListAll"

// SpanSoftDeleteRepo is the span name for soft-deleting a category.
const SpanSoftDeleteRepo = "SoftDelete"

// SpanUpdateRepo is the span name for updating a category.
const SpanUpdateRepo = "Update"

// Span attribute values

// OpCreate is the attribute value for the "create" operation.
const OpCreate = "create"

// OpGetByID is the attribute value for the "get by id" operation.
const OpGetByID = "get_by_id"

// OpGetByName is the attribute value for the "get by name" operation.
const OpGetByName = "get_by_name"

// OpListAll is the attribute value for the "list all" operation.
const OpListAll = "list_all"

// OpSoftDelete is the attribute value for the "soft delete" operation.
const OpSoftDelete = "soft_delete"

// OpUpdate is the attribute value for the "update" operation.
const OpUpdate = "update"

// Status messages (used with span.SetStatus/SetStatus + logs)

// StatusCategoryCreated indicates the category was created successfully.
const StatusCategoryCreated = "category created successfully"

// StatusRetrievedByID indicates a category was retrieved by ID successfully.
const StatusRetrievedByID = "category retrieved by id successfully"

// StatusRetrievedByName indicates a category was retrieved by Name successfully.
const StatusRetrievedByName = "category retrieved by name successfully"

// StatusFetchedAll indicates all categories were retrieved successfully.
const StatusFetchedAll = "all categories retrieved successfully"

// StatusSoftDeleted indicates a category was soft-deleted successfully.
const StatusSoftDeleted = "category soft deleted successfully"

// StatusUpdated indicates a category was updated successfully.
const StatusUpdated = "category updated successfully"

// Error/log messages

// ErrCreateCategoryMsg is the error message used when creation fails.
const ErrCreateCategoryMsg = "error creating category"

// ErrGetCategoryMsg is the error message used when a generic get fails.
const ErrGetCategoryMsg = "error getting category"

// ErrCategoryNotFoundMsg is the error message used when a category is not found.
const ErrCategoryNotFoundMsg = "category not found"
