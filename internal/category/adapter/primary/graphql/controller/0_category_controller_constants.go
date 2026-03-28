// Package controller contains GraphQL-facing controllers for the Category context.
package controller

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for Category GraphQL controllers.
// Format: aion-api.<domain>.<layer>.
const TracerName = "aion-api.category.controller"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreate is the span name for creating a category.
	SpanCreate = "category.controller.create"

	// SpanUpdate is the span name for updating a category.
	SpanUpdate = "category.controller.update"

	// SpanSoftDelete is the span name for soft-deleting a category.
	SpanSoftDelete = "category.controller.soft_delete"

	// SpanListAll is the span name for listing all categories.
	SpanListAll = "category.controller.list_all"

	// SpanGetByID is the span name for retrieving a category by ID.
	SpanGetByID = "category.controller.get_by_id"

	// SpanGetByName is the span name for retrieving a category by name.
	SpanGetByName = "category.controller.get_by_name"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusCreated indicates that a category has been successfully created.
	StatusCreated = "created"

	// StatusUpdated indicates that a category has been successfully updated.
	StatusUpdated = "updated"

	// StatusSoftDeleted indicates that a category has been successfully soft-deleted.
	StatusSoftDeleted = "soft_deleted"

	// StatusFetchedAll indicates that all categories have been successfully fetched.
	StatusFetchedAll = "fetched_all"

	// StatusFetched indicates that a single category has been successfully fetched.
	StatusFetched = "fetched"
)

// =============================================================================
// BUSINESS LOGIC - Error and Log Messages
// =============================================================================

// -----------------------------------------------------------------------------
// Log Messages
// -----------------------------------------------------------------------------

const (
	// MsgCreated is the log message for when a category is created.
	MsgCreated = "category created"

	// MsgUpdated is the log message for when a category is updated.
	MsgUpdated = "category updated"

	// MsgSoftDeleted is the log message for when a category is soft-deleted.
	MsgSoftDeleted = "category soft-deleted"

	// MsgCreateError is the log message for when a create operation fails.
	MsgCreateError = "error creating category"

	// MsgSoftDeleteError is the log message for when a soft-delete operation fails.
	MsgSoftDeleteError = "error soft deleting category"

	// MsgUpdateError is the log message for when an update operation fails.
	MsgUpdateError = "error updating category"
)

// -----------------------------------------------------------------------------
// Sentinel Errors
// Use errors.Is() for type-safe error comparison
// -----------------------------------------------------------------------------

var (
	// ErrUserIDNotFound is the error when the user ID is missing or invalid.
	ErrUserIDNotFound = errors.New("user id not found")

	// ErrCategoryIDNotFound is the error when the category ID is missing or invalid.
	ErrCategoryIDNotFound = errors.New("category id not found")

	// ErrInvalidCategoryID is the error when the category ID cannot be parsed or is invalid.
	ErrInvalidCategoryID = errors.New("invalid category id")

	// ErrCategoryNotFound is the error when a category cannot be found.
	ErrCategoryNotFound = errors.New("category not found")

	// ErrCategoriesNotFound is the error when no categories can be found for a user.
	ErrCategoriesNotFound = errors.New("categories not found")
)
