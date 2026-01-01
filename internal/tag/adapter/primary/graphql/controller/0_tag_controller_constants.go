// Package controller contains GraphQL-facing controllers for the Tag context.
package controller

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for Tag GraphQL controllers.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.tag.controller"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreate is the span name for creating a tag.
	SpanCreate = "tag.controller.create"

	// SpanUpdate is the span name for updating a tag.
	SpanUpdate = "tag.controller.update"

	// SpanGetByName is the span name for retrieving a tag by name.
	SpanGetByName = "tag.controller.get_by_name"

	// SpanGetByCategory is the span name for retrieving tags by category.
	SpanGetByCategory = "tag.controller.get_by_category"

	// SpanListAll is the span name for listing all tags for a user.
	SpanListAll = "tag.controller.list_all"

	// SpanSoftDelete is the span name for soft-deleting a tag.
	SpanSoftDelete = "tag.controller.soft_delete"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusCreated indicates that a tag has been successfully created.
	StatusCreated = "created"

	// StatusUpdated indicates that a tag has been successfully updated.
	StatusUpdated = "updated"

	// StatusFetched indicates that a single tag has been successfully fetched.
	StatusFetched = "fetched"

	// StatusSoftDeleted indicates that a tag has been successfully soft-deleted.
	StatusSoftDeleted = "soft_deleted"
)

// =============================================================================
// BUSINESS LOGIC - Error and Log Messages
// =============================================================================

// -----------------------------------------------------------------------------
// Log Messages
// -----------------------------------------------------------------------------

const (
	// MsgCreated is the log message for when a tag is created.
	MsgCreated = "tag created"

	// MsgCreateError is the log message for when a create operation fails.
	MsgCreateError = "error creating tag"

	// MsgUpdated is the log message for when a tag is updated.
	MsgUpdated = "tag updated"

	// MsgUpdateError is the log message for when an update operation fails.
	MsgUpdateError = "error updating tag"

	// MsgSoftDeleted is the log message for when a tag is soft-deleted.
	MsgSoftDeleted = "tag soft-deleted"

	// MsgSoftDeleteError is the log message for when a soft-delete operation fails.
	MsgSoftDeleteError = "error soft deleting tag"
)

// -----------------------------------------------------------------------------
// Sentinel Errors
// Use errors.Is() for type-safe error comparison
// -----------------------------------------------------------------------------

var (
	// ErrUserIDNotFound is the error when the user ID is missing or invalid.
	ErrUserIDNotFound = errors.New("user id not found")

	// ErrCategoryNotFound is the error when the category ID is missing or invalid.
	ErrCategoryNotFound = errors.New("category id not found")

	// ErrInvalidTagID is the error when the tag ID cannot be parsed or is invalid.
	ErrInvalidTagID = errors.New("invalid tag id")

	// ErrInvalidCategoryID is the error when the category ID cannot be parsed or is invalid.
	ErrInvalidCategoryID = errors.New("invalid category id")

	// ErrTagNotFound is the error when a tag cannot be found.
	ErrTagNotFound = errors.New("tag not found")

	// ErrFailedToListTags is the error when listing tags fails.
	ErrFailedToListTags = errors.New("failed to list tags")
)
