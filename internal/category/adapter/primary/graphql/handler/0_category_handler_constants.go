package handler

const (
	TracerName     = "aionapi.graphql.handler"
	SpanCreate     = "handler.create"
	SpanUpdate     = "handler.update"
	SpanSoftDelete = "handler.soft_delete"
	SpanListAll    = "handler.list_all"
	SpanGetByID    = "handler.get_by_id"
	SpanGetByName  = "handler.get_by_name"

	StatusCreated     = "category_created"
	StatusUpdated     = "category_updated"
	StatusSoftDeleted = "category_soft_deleted"
	StatusFetchedAll  = "all_categories_fetched"
	StatusFetched     = "category_fetched"

	MsgCreated     = "handler created successfully"
	MsgUpdated     = "handler updated successfully"
	MsgSoftDeleted = "handler soft deleted successfully"
	MsgFetchedAll  = "all categories fetched successfully"
	MsgFetched     = "handler fetched successfully"

	ErrInvalidCategoryID = "invalid handler id"
	ErrUserIDNotFound    = "user_id not found in context"

	ErrCategoryNotFound = "handler not found"
)
