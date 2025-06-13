package constants

// TracerUserHandler is the name of the tracer for user operations.
const TracerUserHandler = "AionApi/UserHandler"

// TracerCreateUserHandler is the name of the tracer for user creation operations.
const TracerCreateUserHandler = "CreateUserHandler"

// TracerGetAllUsersHandler is the name of the tracer for user retrieval operations.
const TracerGetAllUsersHandler = "GetAllUsersHandler"

// TracerGetUserByIDHandler is the name of the tracer for user retrieval operations by ID.
const TracerGetUserByIDHandler = "GetUserByIDHandler"

// TracerUpdateUserHandler is the name of the tracer for user update operations.
const TracerUpdateUserHandler = "UpdateUserHandler"

// TracerUpdatePasswordHandler is the name of the tracer for password update operations.
const TracerUpdatePasswordHandler = "UpdatePasswordHandler"

// TracerSoftDeleteUserHandler is the name of the tracer for user soft deletion operations.
const TracerSoftDeleteUserHandler = "SoftDeleteUserHandler"

// UserID is the key used to identify a user in context or request scope.
const UserID = "user_id"

// Username is the key used to identify a user's username in context or request scope.
const Username = "username"

// Email is the key used to identify a user's email in context or request scope.'
const Email = "email"

// ErrorToDecodeUserRequest is returned when decoding a user request fails.
const ErrorToDecodeUserRequest = "error to decode user request"

// ErrorToCreateUser is returned when user creation fails.
const ErrorToCreateUser = "error to create user"

// ErrorToGetUser is returned when getting a user fails.
const ErrorToGetUser = "error to get user"

// ErrorToGetUsers is returned when getting users fails.
const ErrorToGetUsers = "error to get users"

// ErrorToUpdateUser is returned when updating a user fails.
const ErrorToUpdateUser = "error to update user"

// ErrorToSoftDeleteUser is returned when soft deleting a user fails.
const ErrorToSoftDeleteUser = "error to soft delete user"

// ErrorToParseUser is returned when parsing a user fails.
const ErrorToParseUser = "error to parse user"

// ErrorUnauthorizedAccessMissingToken is returned when unauthorized access occurs due to a missing token.
const ErrorUnauthorizedAccessMissingToken = "error unauthorized access missing token" // #nosec G101

// SuccessToCreateUser indicates a successful user creation.
const SuccessToCreateUser = "user created successfully"

// SuccessToGetUser indicates a successful user retrieval.
const SuccessToGetUser = "user get successfully"

// SuccessToGetUsers indicates a successful users retrieval.
const SuccessToGetUsers = "users get successfully"

// SuccessToUpdateUser indicates a successful user update.
const SuccessToUpdateUser = "user updated successfully"

// SuccessToUpdatePassword indicates a successful password update.
const SuccessToUpdatePassword = "password updated successfully"

// SuccessUserSoftDeleted indicates a successful soft deletion of a user.
const SuccessUserSoftDeleted = "user deleted successfully"
