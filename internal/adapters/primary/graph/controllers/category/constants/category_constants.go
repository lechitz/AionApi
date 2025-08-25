package category

const (
	TracerName     = "aionapi.graphql.category"
	SpanCreate     = "category.create"
	SpanUpdate     = "category.update"
	SpanSoftDelete = "category.soft_delete"
	SpanListAll    = "category.list_all"
	SpanGetByID    = "category.get_by_id"
	SpanGetByName  = "category.get_by_name"

	StatusCreated     = "category_created"
	StatusUpdated     = "category_updated"
	StatusSoftDeleted = "category_soft_deleted"
	StatusFetchedAll  = "all_categories_fetched"
	StatusFetched     = "category_fetched"

	MsgCreated     = "category created successfully"
	MsgUpdated     = "category updated successfully"
	MsgSoftDeleted = "category soft deleted successfully"
	MsgFetchedAll  = "all categories fetched successfully"
	MsgFetched     = "category fetched successfully"

	ErrInvalidCategoryID = "invalid category id"
	ErrUserIDNotFound    = "user_id not found in context"

	ErrCategoryNotFound = "category not found"
)
