package constants

const (
	// SpanCreateUser is the name of the span for creating a user.
	SpanCreateUser = "Create"

	// SpanGetSelf is the name of the span for getting a user by ID.
	SpanGetSelf = "GetByID"

	// SpanGetUserByUsername is the name of the span for getting a user by username.
	SpanGetUserByUsername = "FindByUsername"

	// SpanGetUserByEmail is the name of the span for getting a user by email.
	SpanGetUserByEmail = "FindByEmail"

	// SpanGetAllUsers is the name of the span for getting all users.
	SpanGetAllUsers = "ListAll"

	// SpanUpdateUser is the name of the span for updating a user.
	SpanUpdateUser = "Update"

	// SpanUpdateUserPassword is the name of the span for updating a user password.
	SpanUpdateUserPassword = "UdpateUserPassword" // #nosec G101

	// SpanSoftDeleteUser is the name of the span for soft deleting a user.
	SpanSoftDeleteUser = "SoftDelete"

	// StatusValidationFailed is the status for when a validation fails.
	StatusValidationFailed = "validation_failed"

	// StatusDBErrorUsernameOrEmail is the status for when a username check fails.
	StatusDBErrorUsernameOrEmail = "db_error_checking_username_or_email"

	// StatusUsernameExists is the status for when a username already exists.
	StatusUsernameExists = "username_exists"

	// StatusDBErrorEmail is the status for when an email check fails.
	StatusDBErrorEmail = "db_error_checking_email"

	// StatusEmailExists is the status for when an email already exists.
	StatusEmailExists = "email_exists"

	// StatusUsernameAndEmailExists is the status for when a username and email already exists.
	StatusUsernameAndEmailExists = "username_and_email_exists"

	// StatusHashPasswordFailed is the status for when a password hash fails.
	StatusHashPasswordFailed = "hash_password_failed"

	// StatusDBErrorCreateUser is the status for when a user creation fails.
	StatusDBErrorCreateUser = "db_error_create_user"
)
