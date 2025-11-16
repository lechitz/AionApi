// Package controller contains GraphQL-facing controllers for the Tag context.
package controller

// Tracing spans used for OpenTelemetry instrumentation.
const (
	// TracerName is the name of the tracer for Tag GraphQL controllers.
	TracerName = "aionapi.graphql.tag"

	// SpanCreate is the span name for creating a tag.
	SpanCreate = "tag.create"

	// SpanGetByName is the span name for retrieving a tag by name.
	SpanGetByName = "tag.get_by_name"

	// SpanGetByCategory is the span name for retrieving tags by category.
	SpanGetByCategory = "tag.get_by_category"

	// SpanListAll is the span name for listing all tags for a user.
	SpanListAll = "tag.list_all"
)

// Status messages represent the outcome of controller operations.
const (
	// StatusCreated indicates that a tag has been successfully created.
	StatusCreated = "created"

	// StatusFetched indicates that a single tag has been successfully fetched.
	StatusFetched = "fetched"
)

// Log/messages and error strings used in logging and error reporting.
const (
	// MsgCreated is the log message for when a tag is created.
	MsgCreated = "tag created"

	// MsgCreateError is the log message for when a create operation fails.
	MsgCreateError = "error creating tag"

	// ErrUserIDNotFound is the error message when the user ID is missing or invalid.
	ErrUserIDNotFound = "user id not found"

	// ErrCategoryNotFound is the error message when the category ID is missing or invalid.
	ErrCategoryNotFound = "category id not found"

	// ErrInvalidTagID is the error message when the tag ID cannot be parsed or is invalid.
	ErrInvalidTagID = "invalid tag id"

	// ErrInvalidCategoryID is the error message when the category ID cannot be parsed or is invalid.
	ErrInvalidCategoryID = "invalid category id"

	// ErrTagNotFound is the error message when a tag cannot be found.
	ErrTagNotFound = "tag not found"

	// ErrFailedToListTags is the error message when listing tags fails.
	ErrFailedToListTags = "failed to list tags"
)
