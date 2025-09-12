// Package usecase constants contains constants related to handler operations.
package usecase

// Tracing.
const (
	// TracerName is the name of the tracer.
	TracerName = "aionapi.handler"

	// SpanCreateCategory is the name of the span for creating a category.
	SpanCreateCategory = "CreateCategory"

	// SpanGetCategoryByID is the name of the span for getting a category by ID.
	SpanGetCategoryByID = "GetCategoryByID"

	// SpanGetCategoryByName is the name of the span for getting a category by name.
	SpanGetCategoryByName = "GetCategoryByName"

	// SpanListAllCategories is the name of the span for listing all categories.
	SpanListAllCategories = "ListAllCategories"

	// SpanUpdateCategory is the name of the span for updating a category.
	SpanUpdateCategory = "UpdateCategory"

	// SpanSoftDeleteCategory is the name of the span for soft deleting a category.
	SpanSoftDeleteCategory = "SoftDeleteCategory"
)

// Events.
const (
	// EventValidateInput is the name of the event for validating input.
	EventValidateInput = "validate_input"

	// EventCheckUniqueness is the name of the event for checking uniqueness.
	EventCheckUniqueness = "check_uniqueness"

	// EventRepositoryCreate is the name of the event for creating a repository.
	EventRepositoryCreate = "repository_create"

	// EventRepositoryGet is the name of the event for getting a repository.
	EventRepositoryGet = "repository_get"

	// EventRepositoryListAll is the name of the event for listing all repositories.
	EventRepositoryListAll = "repository_list_all"

	// EventRepositoryUpdate is the name of the event for updating a repository.
	EventRepositoryUpdate = "repository_update"

	// EventRepositoryDelete is the name of the event for deleting a repository.
	EventRepositoryDelete = "repository_delete"

	// EventSuccess is the name of the event for a successful operation.
	EventSuccess = "success"
)

// ErrToValidateCategory indicates a validation error in a handler operation.
const ErrToValidateCategory = "handler validation error"

// CategoryAlreadyExists is returned when trying to create a handler that already exists.
const CategoryAlreadyExists = "handler already exists"

// FailedToCreateCategory indicates failure to create a handler.
const FailedToCreateCategory = "failed to create handler"

// FailedToUpdateCategory indicates failure to update a handler.
const FailedToUpdateCategory = "failed to update handler"

// SuccessfullyCreatedCategory is used to format a success message when a handler is created.
const SuccessfullyCreatedCategory = "successfully created handler %s"

// SuccessfullyUpdatedCategory is used to format a success message when a handler is updated.
const SuccessfullyUpdatedCategory = "successfully updated handler %s"

// SuccessfullySoftDeletedCategory is used to format a success message when a handler is soft deleted.
const SuccessfullySoftDeletedCategory = "successfully soft deleted handler %s"

// FailedToGetCategoryByID indicates failure to retrieve a handler by its ID.
const FailedToGetCategoryByID = "failed to get handler by id"

// FailedToGetCategoryByName indicates failure to retrieve a handler by its name.
const FailedToGetCategoryByName = "failed to get handler by name"

// FailedToGetAllCategories indicates failure to retrieve all categories.
const FailedToGetAllCategories = "failed to get all categories"

// FailedToSoftDeleteCategory indicates failure to soft delete a handler.
const FailedToSoftDeleteCategory = "failed to soft delete handler"

// CategoryIDIsRequired indicates the handler ID is required.
const CategoryIDIsRequired = "handler id is required"

// CategoryNameIsRequired indicates the handler name is required.
const CategoryNameIsRequired = "handler name is required"

// CategoryDescriptionIsTooLong indicates the handler description is too long.
const CategoryDescriptionIsTooLong = "handler description cannot exceed 200 characters"

// CategoryColorIsTooLong indicates the handler color is too long.
const CategoryColorIsTooLong = "handler color cannot exceed 7 characters"

// CategoryIconIsTooLong indicates the handler icon is too long.
const CategoryIconIsTooLong = "handler icon cannot exceed 50 characters"
