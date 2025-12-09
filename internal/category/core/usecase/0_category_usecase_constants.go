package usecase

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.category.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreateCategory is the span name for creating a category.
	SpanCreateCategory = "category.create"
	// SpanGetCategoryByID is the span name for getting a category by ID.
	SpanGetCategoryByID = "category.get_by_id"
	// SpanGetCategoryByName is the span name for getting a category by name.
	SpanGetCategoryByName = "category.get_by_name"
	// SpanListAllCategories is the span name for listing all categories.
	SpanListAllCategories = "category.list_all"
	// SpanUpdateCategory is the span name for updating a category.
	SpanUpdateCategory = "category.update"
	// SpanSoftDeleteCategory is the span name for soft-deleting a category.
	SpanSoftDeleteCategory = "category.soft_delete"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// EventValidateInput marks the input-validation step.
	EventValidateInput = "category.input.validate"
	// EventCheckUniqueness marks the uniqueness-check step.
	EventCheckUniqueness = "category.uniqueness.check"
	// EventRepositoryCreate marks the repository create call.
	EventRepositoryCreate = "category.repository.create"
	// EventRepositoryGet marks the repository single-get call.
	EventRepositoryGet = "category.repository.get"
	// EventRepositoryListAll marks the repository list-all call.
	EventRepositoryListAll = "category.repository.list_all"
	// EventRepositoryUpdate marks the repository update call.
	EventRepositoryUpdate = "category.repository.update"
	// EventRepositoryDelete marks the repository delete/soft-delete call.
	EventRepositoryDelete = "category.repository.delete"
	// EventSuccess marks a successful outcome.
	EventSuccess = "category.success"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusCreated indicates a resource was created.
	StatusCreated = "created"
	// StatusRetrievedByID indicates a resource was retrieved by ID.
	StatusRetrievedByID = "retrieved_by_id"
	// StatusRetrievedByName indicates a resource was retrieved by name.
	StatusRetrievedByName = "retrieved_by_name"
	// StatusRetrievedAll indicates resources were retrieved via list-all.
	StatusRetrievedAll = "retrieved_all"
	// StatusUpdated indicates a resource was updated.
	StatusUpdated = "updated"
	// StatusSoftDeleted indicates a resource was deleted (or soft-deleted).
	StatusSoftDeleted = "deleted"
)

// =============================================================================
// BUSINESS LOGIC - Error and Success Messages
// =============================================================================

// Error messages.
const (
	// ErrToValidateCategory indicates a validation error in a category operation.
	ErrToValidateCategory = "category validation error"
	// CategoryAlreadyExists is returned when trying to create a category that already exists.
	CategoryAlreadyExists = "category already exists"
	// FailedToCreateCategory indicates failure to create a category.
	FailedToCreateCategory = "failed to create category"
	// FailedToUpdateCategory indicates failure to update a category.
	FailedToUpdateCategory = "failed to update category"
	// FailedToGetCategoryByID indicates failure to retrieve a category by its ID.
	FailedToGetCategoryByID = "failed to get category by id"
	// FailedToGetCategoryByName indicates failure to retrieve a category by its name.
	FailedToGetCategoryByName = "failed to get category by name"
	// FailedToGetAllCategories indicates failure to retrieve all categories.
	FailedToGetAllCategories = "failed to get all categories"
	// FailedToSoftDeleteCategory indicates failure to soft-delete a category.
	FailedToSoftDeleteCategory = "failed to soft delete category"
)

// Success messages.
const (
	// SuccessfullyCreatedCategory formats a success message when a category is created.
	SuccessfullyCreatedCategory = "successfully created category %s"
	// SuccessfullyUpdatedCategory formats a success message when a category is updated.
	SuccessfullyUpdatedCategory = "successfully updated category %s"
	// SuccessfullySoftDeletedCategory formats a success message when a category is soft-deleted.
	SuccessfullySoftDeletedCategory = "successfully soft deleted category %s"
)

// Validation messages.
const (
	// CategoryIDIsRequired indicates the category ID is required.
	CategoryIDIsRequired = "category id is required"
	// CategoryNameIsRequired indicates the category name is required.
	CategoryNameIsRequired = "category name is required"
	// CategoryDescriptionIsTooLong indicates the category description is too long.
	CategoryDescriptionIsTooLong = "category description cannot exceed 200 characters"
	// CategoryColorIsTooLong indicates the category color is too long.
	CategoryColorIsTooLong = "category color cannot exceed 7 characters"
	// CategoryIconIsTooLong indicates the category icon is too long.
	CategoryIconIsTooLong = "category icon cannot exceed 50 characters"
)
