// Package constants contains constants related to handler operations.
package usecase

// ===== Tracing =====
const (
	TracerName             = "aionapi.handler"
	SpanCreateCategory     = "CreateCategory"
	SpanGetCategoryByID    = "GetCategoryByID"
	SpanGetCategoryByName  = "GetCategoryByName"
	SpanListAllCategories  = "ListAllCategories"
	SpanUpdateCategory     = "UpdateCategory"
	SpanSoftDeleteCategory = "SoftDeleteCategory"
)

// Events (trace)
const (
	EventValidateInput     = "validate_input"
	EventCheckUniqueness   = "check_uniqueness"
	EventRepositoryCreate  = "repository_create"
	EventRepositoryGet     = "repository_get"
	EventRepositoryListAll = "repository_list_all"
	EventRepositoryUpdate  = "repository_update"
	EventRepositoryDelete  = "repository_delete"
	EventSuccess           = "success"
)

// ===== Mensagens/Erros existentes =====

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

// TODO: separar o que foi de pra commonkeys

// CategoryName is the key for handler name in context or responses.
const CategoryName = "name"

// CategoryDescription is the key for handler description in context or responses.
const CategoryDescription = "description"

// CategoryColor is the key for handler color in context or responses.
const CategoryColor = "color_hex"

// CategoryIcon is the key for handler icon in context or responses.
const CategoryIcon = "icon"
