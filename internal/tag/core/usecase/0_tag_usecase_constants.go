package usecase

// TracerName identifies the OpenTelemetry tracer for the Tag usecase package.
const TracerName = "aion.tag.usecase"

// SpanCreateTag is the span name used when creating a new Tag.
const SpanCreateTag = "tag.create"

// SpanGetTagByName is the span name for getting a tag by name.
const SpanGetTagByName = "GetTagByName"

// EventValidateInput marks the moment input validation starts/completes.
const EventValidateInput = "validate_input"

// EventRepositoryCreate marks the step where the repository persistence happens.
const EventRepositoryCreate = "repository_create"

// EventRepositoryGet marks the repository single-get call.
const EventRepositoryGet = "repository_get"

// EventCheckUniqueness marks the uniqueness-check step.
const EventCheckUniqueness = "check_uniqueness"

// TagAlreadyExists is returned when trying to create a tag that already exists.
const TagAlreadyExists = "tag already exists"

// EventSuccess marks the successful completion of the usecase logic.
const EventSuccess = "success"

// ErrToValidateTag is the status description used when input validation fails.
const ErrToValidateTag = "failed_to_validate_tag"

// FailedToCreateTag is the error prefix/message used when the repository fails to persist the tag.
const FailedToCreateTag = "failed_to_create_tag"

// FailedToGetTagByName indicates failure to retrieve a tag by its name.
const FailedToGetTagByName = "failed to get tag by name"

// StatusCreated is the span status description used when the tag is successfully created.
const StatusCreated = "created"

// StatusRetrievedByName indicates a resource was retrieved by name.
const StatusRetrievedByName = "retrieved_by_name"

// SuccessfullyCreatedTag is a structured log message used after a successful creation (expects tag name).
const SuccessfullyCreatedTag = "tag successfully created: %s"

// TagNameIsRequired is the validation message when the tag name is missing.
const TagNameIsRequired = "tag name is required"

// TagDescriptionIsTooLong is the validation message when description exceeds the maximum length.
const TagDescriptionIsTooLong = "tag description is too long"
