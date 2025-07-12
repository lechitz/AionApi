package constants

const (
	// SpanCreateUser is the name of the span for creating a user.
	SpanCreateUser = "CreateUser"

	// StatusValidationFailed is the status for when a validation fails.
	StatusValidationFailed = "validation_failed"

	// StatusDBErrorUsername is the status for when a username check fails.
	StatusDBErrorUsername = "db_error_checking_username"

	// StatusUsernameExists is the status for when a username already exists.
	StatusUsernameExists = "username_exists"

	// StatusDBErrorEmail is the status for when an email check fails.
	StatusDBErrorEmail = "db_error_checking_email"

	// StatusEmailExists is the status for when an email already exists.
	StatusEmailExists = "email_exists"

	// StatusHashPasswordFailed is the status for when a password hash fails.
	StatusHashPasswordFailed = "hash_password_failed"

	// StatusDBErrorCreateUser is the status for when a user creation fails.
	StatusDBErrorCreateUser = "db_error_create_user"

	// StatusSuccess is the status for when a user creation succeeds.
	StatusSuccess = "success"
)
