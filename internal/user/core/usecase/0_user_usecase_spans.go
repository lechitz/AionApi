// Package usecase contains the business logic for the User context.
package usecase

// =============================================================================
// TRACING - OpenTelemetry Instrumentation (Spans)
// =============================================================================

// TracerName is the name of the tracer used for the user use case.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.user.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreateUser is the span name for creating a user.
	SpanCreateUser = "user.create"

	// SpanGetSelf is the span name for getting a user by ID.
	SpanGetSelf = "user.get_by_id"

	// SpanGetUserByUsername is the span name for getting a user by username.
	SpanGetUserByUsername = "user.get_by_username"

	// SpanGetAllUsers is the span name for getting all users.
	SpanGetAllUsers = "user.list_all"

	// SpanUpdateUser is the span name for updating a user.
	SpanUpdateUser = "user.update"

	// SpanUpdateUserPassword is the span name for updating a user password.
	SpanUpdateUserPassword = "user.update_password" // #nosec G101

	// SpanSoftDeleteUser is the span name for soft deleting a user.
	SpanSoftDeleteUser = "user.soft_delete"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusDBErrorUsernameOrEmail is the status for when a username check fails.
	StatusDBErrorUsernameOrEmail = "db_error_checking_username_or_email"

	// StatusHashPasswordFailed is the status for when a password hash fails.
	StatusHashPasswordFailed = "hash_password_failed"

	// StatusDBErrorCreateUser is the status for when a user creation fails.
	StatusDBErrorCreateUser = "db_error_create_user"
)
