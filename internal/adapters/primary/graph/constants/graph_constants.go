// Package constants contains constants used throughout the GraphQL Category/Tag resolvers.
package constants

// --- Tracing: tracer and span names for OpenTelemetry.
const (
	// TracerCategory is the tracer name for Category operations in OpenTelemetry.
	TracerCategory = "aionapi.graphql.category"

	// SpanStartCreateCategory is the span name for creating a category.
	SpanStartCreateCategory = "category.create"

	// SpanStartUpdateCategory is the span name for updating a category.
	SpanStartUpdateCategory = "category.update"

	// SpanStartSoftDeleteCategory is the span name for soft-deleting a category.
	SpanStartSoftDeleteCategory = "category.soft_delete"

	// SpanStartAllGetCategories is the span name for fetching all categories.
	SpanStartAllGetCategories = "category.get_all"

	// SpanStartGetCategoryByID is the span name for fetching a category by ID.
	SpanStartGetCategoryByID = "category.get_by_id"

	// SpanStartGetCategoryByName is the span name for fetching a category by name.
	SpanStartGetCategoryByName = "category.get_by_name"
)

// --- Event names for trace points.
const (
	// EventCreateCategory marks the beginning of a category creation process.
	EventCreateCategory = "category.create.event"

	// EventCategoryCreatedSuccess marks a successful creation of a category.
	EventCategoryCreatedSuccess = "category.created.success"

	// EventUpdateCategory marks the beginning of a category update process.
	EventUpdateCategory = "category.update.event"

	// EventCategoryUpdatedSuccess marks a successful update of a category.
	EventCategoryUpdatedSuccess = "category.updated.success"

	// EventSoftDeleteCategory marks the beginning of a soft delete process for a category.
	EventSoftDeleteCategory = "category.soft_delete.event"

	// EventCategorySoftDeletedSuccess marks a successful soft delete of a category.
	EventCategorySoftDeletedSuccess = "category.soft_deleted.success"

	// EventGetAllCategories marks the fetching of all categories.
	EventGetAllCategories = "category.get_all.event"

	// EventAllCategoriesFetchedSuccess marks successful fetch of all categories.
	EventAllCategoriesFetchedSuccess = "category.fetched_all.success"

	// EventGetCategoryByID marks the fetching of a category by ID.
	EventGetCategoryByID = "category.get_by_id.event"

	// EventGetCategoryByName marks the fetching of a category by name.
	EventGetCategoryByName = "category.get_by_name.event"

	// EventCategoryFetchedSuccess marks a successful fetch of a category.
	EventCategoryFetchedSuccess = "category.fetched.success"
)

// --- Status names for semantic span states.
const (
	// StatusCategoryCreated marks a span as a successful category creation.
	StatusCategoryCreated = "category_created"

	// StatusCategoryUpdated marks a span as a successful category update.
	StatusCategoryUpdated = "category_updated"

	// StatusCategorySoftDeleted marks a span as a successful category soft delete.
	StatusCategorySoftDeleted = "category_soft_deleted"

	// StatusAllCategoriesFetched marks a span as successful fetch of all categories.
	StatusAllCategoriesFetched = "all_categories_fetched"

	// StatusCategoryFetch marks a span as a successful fetch of a category.
	StatusCategoryFetch = "category_fetched"
)

// --- Error and message constants (for logs, responses, etc.).
const (
	// ErrUserIDNotFound is used when user ID is missing from context.
	ErrUserIDNotFound = "user_id not found in context"

	// ErrCategoryCreate is used when there is an error creating a category.
	ErrCategoryCreate = "error creating category"

	// ErrCategoryUpdate is used when there is an error updating a category.
	ErrCategoryUpdate = "error updating category"

	// ErrCategorySoftDelete is used when there is an error soft deleting a category.
	ErrCategorySoftDelete = "error soft deleting category"

	// ErrAllCategoriesNotFound is used when categories cannot be found.
	ErrAllCategoriesNotFound = "categories not found"

	// ErrCategoryNotFound is used when a category cannot be found.
	ErrCategoryNotFound = "category not found"

	// ErrCategoryByNameNotFound is used when a category by name cannot be found.
	ErrCategoryByNameNotFound = "category by name not found"

	// InvalidCategoryID is used when a category ID is invalid.
	InvalidCategoryID = "invalid category id"

	// MsgCategoryCreated is the success message for category creation.
	MsgCategoryCreated = "category created successfully"

	// MsgCategoryUpdated is the success message for category update.
	MsgCategoryUpdated = "category updated successfully"

	// MsgCategorySoftDeleted is the success message for category soft delete.
	MsgCategorySoftDeleted = "category soft deleted successfully"

	// MsgAllCategoriesFetched is the success message for fetching all categories.
	MsgAllCategoriesFetched = "all categories fetched successfully"

	// MsgCategoryFetched is the success message for fetching a category.
	MsgCategoryFetched = "category fetched successfully"
)

// --- Utility for Sprintf struct logging.
const (
	// SprintfStructVerbose is the format string for detailed struct logging.
	SprintfStructVerbose = "%+v"
)
