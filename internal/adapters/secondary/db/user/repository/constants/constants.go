// Package repository holds repository-scoped constants to avoid magic strings.
package repository

const (
	// Tracing / Span names
	TracerUserRepository = "UserRepository"
	SpanCreate           = "Create"
	SpanCheckUniqueness  = "CheckUniqueness"
	SpanGetByID          = "GetByID"
	SpanGetByUsername    = "GetByUsername"
	SpanGetByEmail       = "GetByEmail"
	SpanListAll          = "ListAll"
	SpanUpdate           = "Update"
	SpanSoftDelete       = "SoftDelete"

	// Operation names (attributes)
	OperationCreate          = "create"
	OperationCheckUniqueness = "check_uniqueness"
	OperationGetByID         = "get_by_id"
	OperationGetByUsername   = "get_by_username"
	OperationGetByEmail      = "get_by_email"
	OperationListAll         = "get_all"
	OperationUpdate          = "update"
	OperationSoftDelete      = "soft_delete"

	// Attribute keys (repository-scoped)
	AttrHTTPErrorReason = "http.error_reason"
	LogField            = "field"

	// Postgres constraint identifiers
	PgConstraintUsersUsernameKey = "users_username_key"
	PgConstraintUsersEmailKey    = "users_email_key"

	// Log messages
	LogUniqueViolationOnCreate = "unique constraint violation on create user"
	LogFailedCreateUser        = "failed to create user"
	LogUserCreated             = "user created successfully"

	LogFailedCheckUsername = "failed to check username uniqueness"
	LogFailedCheckEmail    = "failed to check email uniqueness"

	LogFailedGetByID           = "failed to get user by id"
	LogUserRetrievedByID       = "user retrieved by id successfully"
	LogUserNotFoundByUsername  = "user not found by username"
	LogFailedGetByUsername     = "failed to get user by username"
	LogUserRetrievedByUsername = "user retrieved by username successfully"
	LogUserNotFoundByEmail     = "user not found by email"
	LogFailedGetByEmail        = "failed to get user by email"
	LogUserRetrievedByEmail    = "user retrieved by email successfully"

	LogFailedListAll  = "failed to get all users"
	LogUsersRetrieved = "all users retrieved successfully"

	LogFailedUpdateUser = "failed to update user"
	LogUserUpdated      = "user updated successfully"

	LogFailedSoftDelete = "failed to soft delete user"
	LogUserSoftDeleted  = "user soft deleted successfully"

	// Status / span status messages
	StatusValidationDuplicate     = "validation: duplicate"
	SuffixAlreadyExists           = "_already_exists"
	MsgAlreadyInUse               = " is already in use"
	StatusUserCreated             = "user created successfully"
	StatusUniquenessChecked       = "uniqueness checked"
	StatusUserRetrievedByID       = "user retrieved by id successfully"
	StatusUserNotFoundOK          = "user not found (business as usual)"
	StatusUserRetrievedByUsername = "user retrieved by username successfully"
	StatusUserRetrievedByEmail    = "user retrieved by email successfully"
	StatusUsersRetrieved          = "all users retrieved successfully"
	StatusUserUpdated             = "user updated successfully"
	StatusUserSoftDeleted         = "user soft deleted successfully"

	// Column selections (kept here to avoid scattered literals)
	SelectByUsernameColumns = "user_id, name, username, email, password, created_at, updated_at, deleted_at"
	SelectByEmailColumns    = "user_id, name, username, email, password, created_at, updated_at, deleted_at"
	SelectListAllColumns    = "user_id, name, username, email, created_at, updated_at, deleted_at"
)
