// Package repository holds repository-scoped constants to avoid magic strings.
package repository

// Tracing / Span names.
const (
	// TracerUserRepository is the name of the tracer used in the user repository.
	TracerUserRepository = "UserRepository"

	// SpanCreate is the name of the span for creation.
	SpanCreate = "Create"

	// SpanCheckUniqueness is the name of the span for checking uniqueness.
	SpanCheckUniqueness = "CheckUniqueness"

	// SpanGetByID is the name of the span for getting by ID.
	SpanGetByID = "GetByID"

	// SpanGetByUsername is the name of the span for getting by username.
	SpanGetByUsername = "GetByUsername"

	// SpanGetByEmail is the name of the span for getting by email.
	SpanGetByEmail = "GetByEmail"

	// SpanListAll is the name of the span for list all.
	SpanListAll = "ListAll"

	// SpanUpdate is the name of the span for update.
	SpanUpdate = "Update"

	// SpanSoftDelete is the name of the span for soft delete.
	SpanSoftDelete = "SoftDelete"
)

// Operation names (attributes).
const (
	// OperationCreate is the name of the operation for creating a user.
	OperationCreate = "create"

	// OperationCheckUniqueness is the name of the operation for checking uniqueness.
	OperationCheckUniqueness = "check_uniqueness"

	// OperationGetByID is the name of the operation for getting a user by ID.
	OperationGetByID = "get_by_id"

	// OperationGetByUsername is the name of the operation for getting a user by username.
	OperationGetByUsername = "get_by_username"

	// OperationGetByEmail is the name of the operation for getting a user by email.
	OperationGetByEmail = "get_by_email"

	// OperationListAll is the name of the operation for getting all users.
	OperationListAll = "get_all"

	// OperationUpdate is the name of the operation for updating a user.
	OperationUpdate = "update"

	// OperationSoftDelete is the name of the operation for soft deleting a user.
	OperationSoftDelete = "soft_delete"
)

// Attribute keys (repository-scoped).
const (
	// AttrHTTPErrorReason is the name of the attribute for the HTTP error reason.
	AttrHTTPErrorReason = "http.error_reason"

	// LogField is the name of the attribute for the log field.
	LogField = "field"
)

// Postgres constraint identifiers.
const (
	// PgConstraintUsersUsernameKey is the name of the Postgres constraint for username uniqueness.
	PgConstraintUsersUsernameKey = "users_username_key"

	// PgConstraintUsersEmailKey is the name of the Postgres constraint for email uniqueness.
	PgConstraintUsersEmailKey = "users_email_key"
)

// Log messages.
const (
	// LogUniqueViolationOnCreate is the message used when a unique constraint is violated on creation.
	LogUniqueViolationOnCreate = "unique constraint violation on create user"

	// LogFailedCreateUser is the message used when a user creation fails.
	LogFailedCreateUser = "failed to create user"

	// LogUserCreated is the message used when a user is created.
	LogUserCreated = "user created successfully"

	// LogFailedCheckUsername is the message used when a username uniqueness check fails.
	LogFailedCheckUsername = "failed to check username uniqueness"

	// LogFailedCheckEmail is the message used when an email uniqueness check fails.
	LogFailedCheckEmail = "failed to check email uniqueness"

	// LogFailedGetByID is the message used when a user retrieval by ID fails.
	LogFailedGetByID = "failed to get user by id"

	// LogUserRetrievedByID is the message used when ID retrieves a user.
	LogUserRetrievedByID = "user retrieved by id successfully"

	// LogUserNotFoundByUsername is the message used when a username does not find a user.
	LogUserNotFoundByUsername = "user not found by username"

	// LogFailedGetByUsername is the message used when a user retrieval by username fails.
	LogFailedGetByUsername = "failed to get user by username"

	// LogUserRetrievedByUsername is the message used when a username retrieves a user.
	LogUserRetrievedByUsername = "user retrieved by username successfully"

	// LogUserNotFoundByEmail is the message used when a user is not found by email.
	LogUserNotFoundByEmail = "user not found by email"

	// LogFailedGetByEmail is the message used when a user retrieval by email fails.
	LogFailedGetByEmail = "failed to get user by email"

	// LogUserRetrievedByEmail is the message used when a user is retrieved by email.
	LogUserRetrievedByEmail = "user retrieved by email successfully"

	// LogFailedListAll is the message used when a user retrieval by email fails.
	LogFailedListAll = "failed to get all users"

	// LogUsersRetrieved is the message used when all users are retrieved.
	LogUsersRetrieved = "all users retrieved successfully"

	// LogFailedUpdateUser is the message used when a user update fails.
	LogFailedUpdateUser = "failed to update user"

	// LogUserUpdated is the message used when a user is updated.
	LogUserUpdated = "user updated successfully"

	// LogFailedSoftDelete is the message used when a user soft delete fails.
	LogFailedSoftDelete = "failed to soft delete user"

	// LogUserSoftDeleted is the message used when a user is soft deleted.
	LogUserSoftDeleted = "user soft deleted successfully"
)

// Status / span status messages.
const (
	// StatusValidationDuplicate is the status message for a duplicate validation error.
	StatusValidationDuplicate = "validation: duplicate"

	// SuffixAlreadyExists is the suffix for a duplicate error.
	SuffixAlreadyExists = "_already_exists"

	// MsgAlreadyInUse is the message for a duplicate error.
	MsgAlreadyInUse = " is already in use"

	// StatusUserCreated is the status message for a user creation.
	StatusUserCreated = "user created successfully"

	// StatusUniquenessChecked is the status message for a uniqueness check.
	StatusUniquenessChecked = "uniqueness checked"

	// StatusUserRetrievedByID is the status message for a user retrieval by ID.
	StatusUserRetrievedByID = "user retrieved by id successfully"

	// StatusUserNotFoundOK is the status message for a user not found.
	StatusUserNotFoundOK = "user not found (business as usual)"

	// StatusUserRetrievedByUsername is the status message for a user retrieval by username.
	StatusUserRetrievedByUsername = "user retrieved by username successfully"

	// StatusUserRetrievedByEmail is the status message for a user retrieval by email.
	StatusUserRetrievedByEmail = "user retrieved by email successfully"

	// StatusUsersRetrieved is the status message for a user retrieval by email.
	StatusUsersRetrieved = "all users retrieved successfully"

	// StatusUserUpdated is the status message for a user update.
	StatusUserUpdated = "user updated successfully"

	// StatusUserSoftDeleted is the status message for a user soft delete.
	StatusUserSoftDeleted = "user soft deleted successfully"
)

// Column selections (kept here to avoid scattered literals).
const (
	// SelectByUsernameColumns is the columns to select when getting a user by username.
	SelectByUsernameColumns = "user_id, name, username, email, password, roles, created_at, updated_at, deleted_at"

	// SelectByEmailColumns is the columns to select when getting a user by email.
	SelectByEmailColumns = "user_id, name, username, email, password, roles, created_at, updated_at, deleted_at"

	// SelectListAllColumns is the columns to select when listing all users.
	SelectListAllColumns = "user_id, name, username, email, roles, created_at, updated_at, deleted_at"
)
