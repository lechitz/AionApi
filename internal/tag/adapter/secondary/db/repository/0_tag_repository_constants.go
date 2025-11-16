// Package repository defines DB-facing repository logic for Tag.
// This file centralizes all magic strings/constants used by the repository.
package repository

// Tracer & span names (OpenTelemetry)

// TracerName is the tracer name used by the Tag repository.
const TracerName = "TagRepository"

// SpanCreateRepo is the span name for creating a Tag.
const SpanCreateRepo = "Create"

// SpanGetByNameRepo is the span name for fetching a Tag by name.
const SpanGetByNameRepo = "GetByName"

// SpanGetAllRepo is the span name for fetching all Tags for a user.
const SpanGetAllRepo = "GetAll"

// SpanGetByCategoryRepo is the span name for fetching Tags by category.
const SpanGetByCategoryRepo = "GetByCategory"

// Span attribute values

// OpCreate is the attribute value for the "create" operation.
const OpCreate = "create"

// OpGetByName is the attribute value for the "get by name" operation.
const OpGetByName = "get_by_name"

// OpGetByID is the attribute value for the "get by id" operation.
const OpGetByID = "get_by_id"

// OpGetAll is the attribute value for the "get all" operation.
const OpGetAll = "get_all"

// OpGetByCategory is the attribute value for the "get by category" operation.
const OpGetByCategory = "get_by_category"

// Status messages (used with span.SetStatus/SetStatus + logs)

// StatusTagCreated indicates the Tag (handler) was created successfully.
const StatusTagCreated = "Tag created successfully"

// StatusRetrievedByID indicates a Tag was retrieved by ID successfully.
const StatusRetrievedByID = "Tag retrieved by id successfully"

// StatusRetrievedAll indicates all Tags were retrieved successfully.
const StatusRetrievedAll = "Tags retrieved successfully"

// StatusRetrievedByCategory indicates Tags were retrieved by category successfully.
const StatusRetrievedByCategory = "Tags retrieved by category successfully"

// StatusRetrievedByName indicates a Tag was retrieved by Name successfully.
const StatusRetrievedByName = "Tag retrieved by name successfully"

// StatusFetchedAll indicates all categories were retrieved successfully.
const StatusFetchedAll = "all categories retrieved successfully"

// StatusSoftDeleted indicates a Tag was soft-deleted successfully.
const StatusSoftDeleted = "Tag soft deleted successfully"

// StatusUpdated indicates a Tag was updated successfully.
const StatusUpdated = "Tag updated successfully"

// Error/log messages

// ErrCreateTagMsg is the error message used when creation fails.
const ErrCreateTagMsg = "error creating Tag"

// ErrGetTagByIDMsg is the error message used when a generic get fails.
const ErrGetTagByIDMsg = "error getting Tag by ID"

// ErrGetAllTagsMsg is the error message used when getting all Tags fails.
const ErrGetAllTagsMsg = "error getting all Tags for user"

// ErrGetTagsByCategoryMsg is the error message used when getting Tags by category fails.
const ErrGetTagsByCategoryMsg = "error getting Tags by category"

// ErrGetTagByNameMsg is the error message used when a generic get fails.
const ErrGetTagByNameMsg = "error getting Tag by name"

// ErrTagNotFoundMsg is the error message used when a Tag is not found.
const ErrTagNotFoundMsg = "Tag not found"

// ErrNoTagsFoundMsg is the error message used when no Tags are found for a user.
const ErrNoTagsFoundMsg = "no Tags found for user"
