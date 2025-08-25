// Package constants contains constants used throughout the user handler.
package constants

import "errors"

// TracerUserHandler is the tracer name for user handler operations in OpenTelemetry.
const TracerUserHandler = "aionapi.user.handler"

// Span names for OpenTelemetry user handler operations.
const (
	SpanCreateUserHandler       = "user.create"
	SpanGetAllUsersHandler      = "user.get_all"
	SpanGetUserByIDHandler      = "user.get_by_id"
	SpanUpdateUserHandler       = "user.update"
	SpanUpdatePasswordHandler   = "user.update_password"
	SpanSoftDeleteUserHandler   = "user.soft_delete"
	SpanAttrUserPasswordUpdate  = "user.password_update"
	SpanAttrAuthCookieRefreshed = "auth.cookie_refreshed"
)

// Event names for key points within user handler spans.
const (
	EventRequestUserAgentKeyAndIP = "request_user_agent_and_ip"
	EventDecodeRequest            = "decode_request"
	EventUserIDFoundInContext     = "user_id_found_in_context"

	EventUserServiceCreateUser         = "user_service.create_user"
	EventUserServiceGetAllUsers        = "user_service.get_all_users"
	EventUserServiceGetUserByID        = "user_service.get_user_by_id"
	EventUserServiceUpdateUser         = "user_service.update_user"
	EventUserServiceUpdateUserPassword = "user_service.update_user_password"
	EventUserServiceSoftDeleteUser     = "user_service.soft_delete_user"

	EventUserCreatedSuccess         = "user.created.success"
	EventUserUpdatedSuccess         = "user.updated.success"
	EventUserPasswordUpdatedSuccess = "user.password_updated.success"
	EventUserFetchedSuccess         = "user.fetched.success"
	EventUsersFetchedSuccess        = "users.fetched.success"
	EventUserSoftDeletedSuccess     = "user.soft_deleted.success"
)

// Status names for semantic span states.
const (
	StatusUserCreated         = "user_created"
	StatusUsersFetched        = "users_fetched"
	StatusUserFetched         = "user_fetched"
	StatusUserUpdated         = "user_updated"
	StatusUserPasswordUpdated = "user_password_updated"
	StatusUserSoftDeleted     = "user_soft_deleted"
)

// Error messages used in user handler (for response, tracing and logs).
const (
	ErrDecodeGetUserByIDRequest = "decode error on GetByID"

	ErrMissingUserIDParam = "missing user ID parameter"
	ErrInvalidUserIDParam = "invalid user ID"

	ErrUpdateUserPasswordValidation = "validation error on UpdatePassword"

	ErrCreateUser     = "error creating user"
	ErrGetUserByID    = "error getting user by ID"
	ErrGetUsers       = "error getting users"
	ErrUpdateUser     = "error updating user"
	ErrSoftDeleteUser = "error soft deleting user"
)

var ErrNoFieldsToUpdate = errors.New("no fields provided for update")

// Success messages used in user handler.
const (
	MsgUserCreated         = "user created successfully"
	MsgUsersFetched        = "users retrieved successfully"
	MsgUserFetched         = "user retrieved successfully"
	MsgUserUpdated         = "user updated successfully"
	MsgUserPasswordUpdated = "user password updated successfully"
	MsgUserSoftDeleted     = "user soft deleted successfully"
)

// Validation errors and messages for Create handler.
const (
	// ErrCreateUserValidation is the error message for validation errors on Create.
	ErrCreateUserValidation = "validation error on Create"
)
