package usecase

import "errors"

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
	// EventInvalidateCache marks cache invalidation operations.
	EventInvalidateCache = "category.cache.invalidate"
	// EventSaveToCache marks cache save operations.
	EventSaveToCache = "category.cache.save"
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
	// CategoryIconInvalid indicates the category icon is invalid.
	CategoryIconInvalid = "category icon must be a valid svg key (ex: health/brain.svg)"
)

// Cache operation messages.
const (
	// WarnFailedToInvalidateCategoryCache indicates failure to invalidate single category cache.
	WarnFailedToInvalidateCategoryCache = "failed to invalidate category cache after soft delete"
	// WarnFailedToInvalidateCategoryListCache indicates failure to invalidate category list cache.
	WarnFailedToInvalidateCategoryListCache = "failed to invalidate category list cache after creating category"
	// WarnFailedToDeleteCategoryCache indicates failure to delete category cache after update.
	WarnFailedToDeleteCategoryCache = "failed to delete category cache after update"
	// WarnFailedToDeleteCategoryByNameCache indicates failure to delete category-by-name cache.
	WarnFailedToDeleteCategoryByNameCache = "failed to delete category-by-name cache after update"
	// WarnFailedToSaveCategoryToCache indicates failure to save category to cache by ID.
	WarnFailedToSaveCategoryToCache = "failed to save category to cache"
	// WarnFailedToSaveCategoryToCacheByName indicates failure to save category to cache by name.
	WarnFailedToSaveCategoryToCacheByName = "failed to save category to cache by name"
	// WarnFailedToSaveCategoryToCacheByID indicates failure to save category to cache by ID.
	WarnFailedToSaveCategoryToCacheByID = "failed to save category to cache by ID"
	// WarnFailedToSaveCategoryListToCache indicates failure to save category list to cache.
	WarnFailedToSaveCategoryListToCache = "failed to save category list to cache"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrValidateCategory is a sentinel error for category validation failures.
	ErrValidateCategory = errors.New(ErrToValidateCategory)

	// ErrCategoryAlreadyExists is a sentinel error when a category already exists.
	ErrCategoryAlreadyExists = errors.New(CategoryAlreadyExists)

	// ErrCreateCategory is a sentinel error for category creation failures.
	ErrCreateCategory = errors.New(FailedToCreateCategory)

	// ErrUpdateCategory is a sentinel error for category update failures.
	ErrUpdateCategory = errors.New(FailedToUpdateCategory)

	// ErrGetCategoryByID is a sentinel error for retrieving category by ID.
	ErrGetCategoryByID = errors.New(FailedToGetCategoryByID)

	// ErrGetCategoryByName is a sentinel error for retrieving category by name.
	ErrGetCategoryByName = errors.New(FailedToGetCategoryByName)

	// ErrGetAllCategories is a sentinel error for listing all categories.
	ErrGetAllCategories = errors.New(FailedToGetAllCategories)

	// ErrSoftDeleteCategory is a sentinel error for soft delete failures.
	ErrSoftDeleteCategory = errors.New(FailedToSoftDeleteCategory)

	// ErrCategoryIDRequired is a sentinel error when category ID is missing.
	ErrCategoryIDRequired = errors.New(CategoryIDIsRequired)

	// ErrCategoryNameRequired is a sentinel error when category name is missing.
	ErrCategoryNameRequired = errors.New(CategoryNameIsRequired)

	// ErrCategoryDescriptionTooLong is a sentinel error when description exceeds limit.
	ErrCategoryDescriptionTooLong = errors.New(CategoryDescriptionIsTooLong)

	// ErrCategoryColorTooLong is a sentinel error when color exceeds limit.
	ErrCategoryColorTooLong = errors.New(CategoryColorIsTooLong)

	// ErrCategoryIconInvalid is a sentinel error when icon is invalid.
	ErrCategoryIconInvalid = errors.New(CategoryIconInvalid)
)
