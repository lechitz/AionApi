// Package controller contains GraphQL-facing controllers for the Category context.
package controller

// Tracing spans used for OpenTelemetry instrumentation.
const (
	// TracerName is the name of the tracer for Category GraphQL controllers.
	TracerName = "aionapi.graphql.category"

	// SpanCreate is the span name for creating a category.
	SpanCreate = "category.create"

	// SpanUpdate is the span name for updating a category.
	SpanUpdate = "category.update"

	// SpanSoftDelete is the span name for soft-deleting a category.
	SpanSoftDelete = "category.soft_delete"

	// SpanListAll is the span name for listing all categories.
	SpanListAll = "category.list_all"

	// SpanGetByID is the span name for retrieving a category by ID.
	SpanGetByID = "category.get_by_id"

	// SpanGetByName is the span name for retrieving a category by name.
	SpanGetByName = "category.get_by_name"
)

// Status messages represent the outcome of controller operations.
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

// Log/messages and error strings used in logging and error reporting.
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

	// ErrUserIDNotFound is the error message when the user ID is missing or invalid.
	ErrUserIDNotFound = "user id not found"

	// ErrCategoryIDNotFound is the error message when the category ID is missing or invalid.
	ErrCategoryIDNotFound = "category id not found"

	// ErrInvalidCategoryID is the error message when the category ID cannot be parsed or is invalid.
	ErrInvalidCategoryID = "invalid category id"

	// ErrCategoryNotFound is the error message when a category cannot be found.
	ErrCategoryNotFound = "category not found"

	// ErrCategoriesNotFound is the error message when no categories can be found for a user.
	ErrCategoriesNotFound = "categories not found"
)
