// Package handler contains the GraphQL handlers for the category service.
package handler

const (
	// TracerName is the name of the tracer.
	TracerName = "aionapi.graphql.handler"

	// SpanCreate is the name of the span for create.
	SpanCreate = "handler.create"

	// SpanUpdate is the name of the span for update.
	SpanUpdate = "handler.update"

	// SpanSoftDelete is the name of the span for soft delete.
	SpanSoftDelete = "handler.soft_delete"

	// SpanListAll is the name of the span for list all.
	SpanListAll = "handler.list_all"

	// SpanGetByID is the name of the span for get by id.
	SpanGetByID = "handler.get_by_id"

	// SpanGetByName is the name of the span for get by name.
	SpanGetByName = "handler.get_by_name"

	// StatusCreated is the status for when a handler is created.
	StatusCreated = "category_created"

	// StatusUpdated is the status for when a handler is updated.
	StatusUpdated = "category_updated"

	// StatusSoftDeleted is the status for when a handler is soft deleted.
	StatusSoftDeleted = "category_soft_deleted"

	// StatusFetchedAll is the status for when all handlers are fetched.
	StatusFetchedAll = "all_categories_fetched"

	// StatusFetched is the status for when a handler is fetched.
	StatusFetched = "category_fetched"

	// MsgCreated is the message for when a handler is created.
	MsgCreated = "handler created successfully"

	// MsgUpdated is the message for when a handler is updated.
	MsgUpdated = "handler updated successfully"

	// MsgSoftDeleted is the message for when a handler is soft deleted.
	MsgSoftDeleted = "handler soft deleted successfully"

	// MsgFetchedAll is the message for when all handlers are fetched.
	MsgFetchedAll = "all categories fetched successfully"

	// MsgFetched is the message for when a handler is fetched.
	MsgFetched = "handler fetched successfully"

	// ErrInvalidCategoryID is the error message for when an invalid category ID is provided.
	ErrInvalidCategoryID = "invalid handler id"

	// ErrUserIDNotFound is the error message for when the user ID is not found in the context.
	ErrUserIDNotFound = "user_id not found in context"

	// ErrCategoryNotFound is the error message for when a handler is not found.
	ErrCategoryNotFound = "handler not found"
)
