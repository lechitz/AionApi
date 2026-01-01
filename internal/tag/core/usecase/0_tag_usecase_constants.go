package usecase

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName identifies the OpenTelemetry tracer for the Tag usecase package.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.tag.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreateTag is the span name used when creating a new Tag.
	SpanCreateTag = "tag.create"

	// SpanUpdateTag is the span name used when updating a Tag.
	SpanUpdateTag = "tag.update"

	// SpanGetTagByName is the span name for getting a tag by name.
	SpanGetTagByName = "tag.get_by_name"

	// SpanGetByCategory is the span name for retrieving tags by category.
	SpanGetByCategory = "tag.get_by_category"

	// SpanGetAll is the span name for listing all tags for a user.
	SpanGetAll = "tag.list_all"

	// SpanSoftDeleteTag is the span name for soft-deleting a tag.
	SpanSoftDeleteTag = "tag.soft_delete"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// EventValidateInput marks the moment input validation starts/completes.
	EventValidateInput = "tag.input.validate"

	// EventRepositoryCreate marks the step where the repository persistence happens.
	EventRepositoryCreate = "tag.repository.create"

	// EventRepositoryUpdate marks the step where the repository update happens.
	EventRepositoryUpdate = "tag.repository.update"

	// EventRepositoryGet marks the repository single-get call.
	EventRepositoryGet = "tag.repository.get"

	// EventRepositoryListAll marks the repository list-all call.
	EventRepositoryListAll = "tag.repository.list_all"

	// EventCheckUniqueness marks the uniqueness-check step.
	EventCheckUniqueness = "tag.uniqueness.check"

	// EventInvalidateCache marks the cache invalidation step.
	EventInvalidateCache = "tag.cache.invalidate"

	// EventSuccess marks the successful completion of the usecase logic.
	EventSuccess = "tag.success"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusCreated is the span status description used when the tag is successfully created.
	StatusCreated = "created"

	// StatusUpdated is the span status description used when the tag is successfully updated.
	StatusUpdated = "updated"

	// StatusRetrievedByName indicates a resource was retrieved by name.
	StatusRetrievedByName = "retrieved_by_name"

	// StatusListedAll indicates all tags were listed successfully.
	StatusListedAll = "listed_all"

	// StatusSoftDeleted indicates a resource was soft-deleted.
	StatusSoftDeleted = "deleted"
)

// =============================================================================
// BUSINESS LOGIC - Error and Success Messages
// =============================================================================

// Error messages.
const (
	// ErrToValidateTag is the status description used when input validation fails.
	ErrToValidateTag = "failed_to_validate_tag"

	// FailedToCreateTag is the error prefix/message used when the repository fails to persist the tag.
	FailedToCreateTag = "failed_to_create_tag"

	// FailedToUpdateTag is the error prefix/message used when the repository fails to update the tag.
	FailedToUpdateTag = "failed_to_update_tag"

	// FailedToGetTagByName indicates failure to retrieve a tag by its name.
	FailedToGetTagByName = "failed to get tag by name"

	// FailedToListTags indicates failure to list all tags for a user.
	FailedToListTags = "failed_to_list_tags"

	// ErrFailedToListTags is the error message used when listing tags fails.
	ErrFailedToListTags = "failed to list tags"

	// TagAlreadyExists is returned when trying to create a tag that already exists.
	TagAlreadyExists = "tag already exists"

	// FailedToSoftDeleteTag indicates failure to soft-delete a tag.
	FailedToSoftDeleteTag = "failed to soft delete tag"

	// FailedToInvalidateTagListCache indicates failure to invalidate tag list cache.
	FailedToInvalidateTagListCache = "failed to invalidate tag list cache"

	// FailedToInvalidateTagsByCategoryCache indicates failure to invalidate tags by category cache.
	FailedToInvalidateTagsByCategoryCache = "failed to invalidate tags by category cache"
)

// Success messages.
const (
	// SuccessfullyCreatedTag is a structured log message used after a successful creation.
	SuccessfullyCreatedTag = "tag successfully created: %s"

	// SuccessfullyUpdatedTag is a structured log message used after a successful update.
	SuccessfullyUpdatedTag = "tag successfully updated"

	// SuccessfullySoftDeletedTag formats a success message when a tag is soft-deleted.
	SuccessfullySoftDeletedTag = "successfully soft deleted tag %s"
)

// Validation messages.
const (
	// UserIDIsRequired is the validation message when the user ID is missing.
	UserIDIsRequired = "user ID is required"

	// TagNameIsRequired is the validation message when the tag name is missing.
	TagNameIsRequired = "tag name is required"

	// TagDescriptionIsTooLong is the validation message when description exceeds the maximum length.
	TagDescriptionIsTooLong = "tag description is too long"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrValidateTag is a sentinel error for tag validation failures.
	ErrValidateTag = errors.New(ErrToValidateTag)

	// ErrCreateTag is a sentinel error for tag creation failures.
	ErrCreateTag = errors.New(FailedToCreateTag)

	// ErrUpdateTag is a sentinel error for tag update failures.
	ErrUpdateTag = errors.New(FailedToUpdateTag)

	// ErrGetTagByName is a sentinel error for retrieving tag by name.
	ErrGetTagByName = errors.New(FailedToGetTagByName)

	// ErrListTags is a sentinel error for listing tags.
	ErrListTags = errors.New(ErrFailedToListTags)

	// ErrTagAlreadyExists is a sentinel error when a tag already exists.
	ErrTagAlreadyExists = errors.New(TagAlreadyExists)

	// ErrSoftDeleteTag is a sentinel error for soft delete failures.
	ErrSoftDeleteTag = errors.New(FailedToSoftDeleteTag)

	// ErrInvalidateTagListCache is a sentinel error for cache invalidation failures.
	ErrInvalidateTagListCache = errors.New(FailedToInvalidateTagListCache)

	// ErrInvalidateTagsByCategoryCache is a sentinel error for cache invalidation failures.
	ErrInvalidateTagsByCategoryCache = errors.New(FailedToInvalidateTagsByCategoryCache)

	// ErrUserIDRequired is a sentinel error when user ID is missing.
	ErrUserIDRequired = errors.New(UserIDIsRequired)

	// ErrTagNameRequired is a sentinel error when tag name is missing.
	ErrTagNameRequired = errors.New(TagNameIsRequired)

	// ErrTagDescriptionTooLong is a sentinel error when description exceeds limit.
	ErrTagDescriptionTooLong = errors.New(TagDescriptionIsTooLong)
)
