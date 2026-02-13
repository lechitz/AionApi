// Package repository holds repository-scoped constants to avoid magic strings.
package repository

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer used in the user repository.
// Format: aionapi.<domain>.<layer>.
const TracerName = "aionapi.user.repository"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreate is the span name for creation.
	SpanCreate = "user.repository.create"

	// SpanCheckUniqueness is the span name for checking uniqueness.
	SpanCheckUniqueness = "user.repository.check_uniqueness"

	// SpanGetByID is the span name for getting by ID.
	SpanGetByID = "user.repository.get_by_id"

	// SpanGetByUsername is the span name for getting by username.
	SpanGetByUsername = "user.repository.get_by_username"

	// SpanListAll is the span name for list all.
	SpanListAll = "user.repository.list_all"

	// SpanUpdate is the span name for update.
	SpanUpdate = "user.repository.update"

	// SpanSoftDelete is the span name for soft delete.
	SpanSoftDelete = "user.repository.soft_delete"
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// OperationCreate is the attribute value for creating a user.
	OperationCreate = "create"

	// OperationCheckUniqueness is the attribute value for checking uniqueness.
	OperationCheckUniqueness = "check_uniqueness"

	// OperationGetByID is the attribute value for getting a user by ID.
	OperationGetByID = "get_by_id"

	// OperationGetByUsername is the attribute value for getting a user by username.
	OperationGetByUsername = "get_by_username"

	// OperationListAll is the attribute value for getting all users.
	OperationListAll = "get_all"

	// OperationUpdate is the attribute value for updating a user.
	OperationUpdate = "update"

	// OperationSoftDelete is the attribute value for soft deleting a user.
	OperationSoftDelete = "soft_delete"
)

// Attribute keys (repository-scoped).
const (
	// AttrHTTPErrorReason is the attribute key for the HTTP error reason.
	AttrHTTPErrorReason = "http.error_reason"

	// LogField is the attribute key for the log field.
	LogField = "field"
)

// =============================================================================
// BUSINESS LOGIC - Postgres Constraints
// =============================================================================

const (
	// PgConstraintUsersUsernameKey is the Postgres constraint for username uniqueness.
	PgConstraintUsersUsernameKey = "users_username_key"

	// PgConstraintUsersEmailKey is the Postgres constraint for email uniqueness.
	PgConstraintUsersEmailKey = "users_email_key"
)

// =============================================================================
// BUSINESS LOGIC - Log Messages
// =============================================================================

const (
	// LogUniqueViolationOnCreate is logged when a unique constraint is violated on creation.
	LogUniqueViolationOnCreate = "unique constraint violation on create user"

	// LogFailedCreateUser is logged when a user creation fails.
	LogFailedCreateUser = "failed to create user"

	// LogUserCreated is logged when a user is created.
	LogUserCreated = "user created successfully"

	// LogFailedCheckUsername is logged when a username uniqueness check fails.
	LogFailedCheckUsername = "failed to check username uniqueness"

	// LogFailedCheckEmail is logged when an email uniqueness check fails.
	LogFailedCheckEmail = "failed to check email uniqueness"

	// LogFailedGetByID is logged when a user retrieval by ID fails.
	LogFailedGetByID = "failed to get user by id"

	// LogUserRetrievedByID is logged when ID retrieves a user.
	LogUserRetrievedByID = "user retrieved by id successfully"

	// LogUserNotFoundByUsername is logged when a username does not find a user.
	LogUserNotFoundByUsername = "user not found by username"

	// LogFailedGetByUsername is logged when a user retrieval by username fails.
	LogFailedGetByUsername = "failed to get user by username"

	// LogUserRetrievedByUsername is logged when a username retrieves a user.
	LogUserRetrievedByUsername = "user retrieved by username successfully"

	// LogFailedListAll is logged when listing all users fails.
	LogFailedListAll = "failed to get all users"

	// LogUsersRetrieved is logged when all users are retrieved.
	LogUsersRetrieved = "all users retrieved successfully"

	// LogFailedUpdateUser is logged when a user update fails.
	LogFailedUpdateUser = "failed to update user"

	// LogUserUpdated is logged when a user is updated.
	LogUserUpdated = "user updated successfully"

	// LogFailedSoftDelete is logged when a user soft delete fails.
	LogFailedSoftDelete = "failed to soft delete user"

	// LogUserSoftDeleted is logged when a user is soft deleted.
	LogUserSoftDeleted = "user soft deleted successfully"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusValidationDuplicate is the status for a duplicate validation error.
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
	// Includes optional profile fields because auth login caches this user payload.
	SelectByUsernameColumns = "user_id, name, username, email, password, locale, timezone, location, bio, avatar_url, created_at, updated_at, deleted_at"

	// SelectListAllColumns is the columns to select when listing all users.
	SelectListAllColumns = "user_id, name, username, email, created_at, updated_at, deleted_at"
)
