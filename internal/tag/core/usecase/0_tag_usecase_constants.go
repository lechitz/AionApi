package usecase

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

	// SpanGetTagByName is the span name for getting a tag by name.
	SpanGetTagByName = "tag.get_by_name"

	// SpanGetByCategory is the span name for retrieving tags by category.
	SpanGetByCategory = "tag.get_by_category"

	// SpanGetAll is the span name for listing all tags for a user.
	SpanGetAll = "tag.list_all"
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

	// EventRepositoryGet marks the repository single-get call.
	EventRepositoryGet = "tag.repository.get"

	// EventRepositoryListAll marks the repository list-all call.
	EventRepositoryListAll = "tag.repository.list_all"

	// EventCheckUniqueness marks the uniqueness-check step.
	EventCheckUniqueness = "tag.uniqueness.check"

	// EventSuccess marks the successful completion of the usecase logic.
	EventSuccess = "tag.success"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusCreated is the span status description used when the tag is successfully created.
	StatusCreated = "created"

	// StatusRetrievedByName indicates a resource was retrieved by name.
	StatusRetrievedByName = "retrieved_by_name"

	// StatusListedAll indicates all tags were listed successfully.
	StatusListedAll = "listed_all"
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

	// FailedToGetTagByName indicates failure to retrieve a tag by its name.
	FailedToGetTagByName = "failed to get tag by name"

	// FailedToListTags indicates failure to list all tags for a user.
	FailedToListTags = "failed_to_list_tags"

	// ErrFailedToListTags is the error message used when listing tags fails.
	ErrFailedToListTags = "failed to list tags"

	// TagAlreadyExists is returned when trying to create a tag that already exists.
	TagAlreadyExists = "tag already exists"
)

// Success messages.
const (
	// SuccessfullyCreatedTag is a structured log message used after a successful creation.
	SuccessfullyCreatedTag = "tag successfully created: %s"
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
