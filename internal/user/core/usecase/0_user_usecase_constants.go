// Package usecase constants contains constants related to user operations.
package usecase

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
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

	// SpanGetUserStats is the span name for getting user statistics.
	SpanGetUserStats = "user.get_stats"

	// SpanUpdateUser is the span name for updating a user.
	SpanUpdateUser = "user.update"

	// SpanUpdateUserPassword is the span name for updating a user password.
	SpanUpdateUserPassword = "user.update_password" // #nosec G101

	// SpanSoftDeleteUser is the span name for soft deleting a user.
	SpanSoftDeleteUser = "user.soft_delete"

	// SpanUploadAvatar is the span name for uploading/processing avatar images.
	SpanUploadAvatar = "user.upload_avatar"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// SpanEventCheckCache is the event name for checking cache.
	SpanEventCheckCache = "CheckCache"

	// SpanEventCacheHit is the event name when cache hit occurs.
	SpanEventCacheHit = "CacheHit"

	// SpanEventCacheMiss is the event name when cache miss occurs.
	SpanEventCacheMiss = "CacheMiss"

	// SpanEventSaveToCache is the event name for saving to cache.
	SpanEventSaveToCache = "SaveToCache"

	// SpanEventInvalidateCache is the event name for invalidating cache.
	SpanEventInvalidateCache = "InvalidateCache"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusDBErrorUsernameOrEmail is the status for when a username check fails.
	StatusDBErrorUsernameOrEmail = "db_error_checking_username_or_email"

	// StatusUsernameOrEmailInUse is the status for when username or email already exists.
	StatusUsernameOrEmailInUse = "username_or_email_in_use"

	// StatusHashPasswordFailed is the status for when a password hash fails.
	StatusHashPasswordFailed = "hash_password_failed"

	// StatusDBErrorCreateUser is the status for when a user creation fails.
	StatusDBErrorCreateUser = "db_error_create_user"
)

// =============================================================================
// BUSINESS LOGIC - Roles
// =============================================================================

// UserRoles is the role of a user.
const UserRoles = "user"

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrorToHashPassword indicates an error while hashing a password.
	// #nosec G101: This constant does not leak a real secret, just an error message.
	ErrorToHashPassword = "error hashing password"

	// ErrorToCreateUser indicates an error when creating a user.
	ErrorToCreateUser = "error to create user"

	// ErrorToCompareHashAndPassword indicates a password hash comparison failure.
	ErrorToCompareHashAndPassword = "error to compare hash and password"

	// ErrorToCreateToken indicates an error when creating a token.
	ErrorToCreateToken = "error to create token"

	// ErrorToGetSelf indicates an error when fetching a user by ID.
	ErrorToGetSelf = "error to get user by id"

	// ErrorNoFieldsToUpdate indicates there were no fields to update in the user.
	ErrorNoFieldsToUpdate = "no fields to update"

	// ErrorToUpdatePassword indicates an error when updating the user password.
	ErrorToUpdatePassword = "error to update password"

	// ErrorToUpdateUser indicates an error when updating the user.
	ErrorToUpdateUser = "error to update user"

	// ErrorToGetUserByUsername indicates an error when fetching a user by username.
	ErrorToGetUserByUsername = "error to get user by username"

	// ErrorToSoftDeleteUser indicates an error when performing a soft delete on a user.
	ErrorToSoftDeleteUser = "error to soft delete user"

	// ErrorToGetUserStats indicates an error when getting user statistics.
	ErrorToGetUserStats = "error getting user stats"
)

// =============================================================================
// BUSINESS LOGIC - Success Messages
// =============================================================================

const (
	// SuccessUserCreated indicates that the user was created successfully.
	SuccessUserCreated = "user created successfully"

	// SuccessUserRetrieved indicates a user was successfully retrieved.
	SuccessUserRetrieved = "user retrieved successfully"

	// SuccessPasswordUpdated indicates the password was updated successfully.
	SuccessPasswordUpdated = "password updated successfully"

	// SuccessUserUpdated indicates the user was updated successfully.
	SuccessUserUpdated = "user updated successfully"

	// SuccessUserSoftDeleted indicates a user was softly deleted successfully.
	SuccessUserSoftDeleted = "user soft deleted successfully"

	// SuccessUserStatsRetrieved indicates user stats were retrieved successfully.
	SuccessUserStatsRetrieved = "user stats retrieved successfully"
)

// =============================================================================
// LOGGING - Info Messages
// =============================================================================

const (
	// InfoUserRetrievedFromCache is an info message when user is retrieved from cache.
	InfoUserRetrievedFromCache = "user retrieved from cache"

	// InfoGettingUserStats is an info message when getting user stats.
	InfoGettingUserStats = "getting user stats"
)

// =============================================================================
// LOGGING - Warning Messages
// =============================================================================

const (
	// WarnUsernameOrEmailInUse is a warning message when username or email is already in use.
	WarnUsernameOrEmailInUse = "username or email already in use"

	// WarnFailedToSaveUserToCache is a warning message when saving user to cache fails.
	WarnFailedToSaveUserToCache = "failed to save user to cache after creation"

	// WarnFailedToSaveUserToCacheGeneric is a warning message when saving user to cache fails (generic).
	WarnFailedToSaveUserToCacheGeneric = "failed to save user to cache"

	// WarnFailedToInvalidateUserCache is a warning message when invalidating user cache fails.
	WarnFailedToInvalidateUserCache = "failed to invalidate user cache after update"

	// WarnFailedToInvalidateUserCacheAfterDelete is a warning message when invalidating user cache fails after soft delete.
	WarnFailedToInvalidateUserCacheAfterDelete = "failed to invalidate user cache after soft delete"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrHashPassword is a sentinel error for password hashing failures.
	ErrHashPassword = errors.New(ErrorToHashPassword)

	// ErrCreateUser is a sentinel error for user creation failures.
	ErrCreateUser = errors.New(ErrorToCreateUser)

	// ErrGetSelf is a sentinel error for retrieving user by ID.
	ErrGetSelf = errors.New(ErrorToGetSelf)

	// ErrNoFieldsToUpdate is a sentinel error when no fields are provided for update.
	ErrNoFieldsToUpdate = errors.New(ErrorNoFieldsToUpdate)

	// ErrUpdatePassword is a sentinel error for password update failures.
	ErrUpdatePassword = errors.New(ErrorToUpdatePassword)

	// ErrUpdateUser is a sentinel error for user update failures.
	ErrUpdateUser = errors.New(ErrorToUpdateUser)

	// ErrCompareHashAndPassword is a sentinel error for password comparison failures.
	ErrCompareHashAndPassword = errors.New(ErrorToCompareHashAndPassword)

	// ErrCreateToken is a sentinel error for token creation failures.
	ErrCreateToken = errors.New(ErrorToCreateToken)

	// ErrGetUserByUsername is a sentinel error for retrieving user by username.
	ErrGetUserByUsername = errors.New(ErrorToGetUserByUsername)

	// ErrSoftDeleteUser is a sentinel error for soft delete failures.
	ErrSoftDeleteUser = errors.New(ErrorToSoftDeleteUser)
)
