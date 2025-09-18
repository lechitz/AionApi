package usecase

// Tracing.
const (
	// TracerName is the name of the tracer.
	TracerName = "aionapi.handler"

	// SpanCreateCategory is the span name for creating a category.
	SpanCreateCategory = "CreateCategory"
	// SpanGetCategoryByID is the span name for getting a category by ID.
	SpanGetCategoryByID = "GetCategoryByID"
	// SpanGetCategoryByName is the span name for getting a category by name.
	SpanGetCategoryByName = "GetCategoryByName"
	// SpanListAllCategories is the span name for listing all categories.
	SpanListAllCategories = "ListAllCategories"
	// SpanUpdateCategory is the span name for updating a category.
	SpanUpdateCategory = "UpdateCategory"
	// SpanSoftDeleteCategory is the span name for soft-deleting a category.
	SpanSoftDeleteCategory = "SoftDeleteCategory"
)

// Events.
const (
	// EventValidateInput marks the input-validation step.
	EventValidateInput = "validate_input"
	// EventCheckUniqueness marks the uniqueness-check step.
	EventCheckUniqueness = "check_uniqueness"
	// EventRepositoryCreate marks the repository create call.
	EventRepositoryCreate = "repository_create"
	// EventRepositoryGet marks the repository single-get call.
	EventRepositoryGet = "repository_get"
	// EventRepositoryListAll marks the repository list-all call.
	EventRepositoryListAll = "repository_list_all"
	// EventRepositoryUpdate marks the repository update call.
	EventRepositoryUpdate = "repository_update"
	// EventRepositoryDelete marks the repository delete/soft-delete call.
	EventRepositoryDelete = "repository_delete"
	// EventSuccess marks a successful outcome.
	EventSuccess = "success"
)

// Span status messages (used with span.SetStatus).
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

// Domain/user-facing error/status messages.
const (
	// ErrToValidateCategory indicates a validation error in a category operation.
	ErrToValidateCategory = "category validation error"
	// CategoryAlreadyExists is returned when trying to create a category that already exists.
	CategoryAlreadyExists = "category already exists"
	// FailedToCreateCategory indicates failure to create a category.
	FailedToCreateCategory = "failed to create category"
	// FailedToUpdateCategory indicates failure to update a category.
	FailedToUpdateCategory = "failed to update category"
	// SuccessfullyCreatedCategory formats a success message when a category is created.
	SuccessfullyCreatedCategory = "successfully created category %s"
	// SuccessfullyUpdatedCategory formats a success message when a category is updated.
	SuccessfullyUpdatedCategory = "successfully updated category %s"
	// SuccessfullySoftDeletedCategory formats a success message when a category is soft-deleted.
	SuccessfullySoftDeletedCategory = "successfully soft deleted category %s"
	// FailedToGetCategoryByID indicates failure to retrieve a category by its ID.
	FailedToGetCategoryByID = "failed to get category by id"
	// FailedToGetCategoryByName indicates failure to retrieve a category by its name.
	FailedToGetCategoryByName = "failed to get category by name"
	// FailedToGetAllCategories indicates failure to retrieve all categories.
	FailedToGetAllCategories = "failed to get all categories"
	// FailedToSoftDeleteCategory indicates failure to soft-delete a category.
	FailedToSoftDeleteCategory = "failed to soft delete category"
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
