// Package usecase contains business logic for the Admin context.
package usecase

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer used for the admin use case.
// Format: aion-api.<domain>.<layer> .
const TracerName = "aion-api.admin.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanHealthCheck is the span name for health check operations.
	SpanHealthCheck = "admin.health_check"

	// SpanGetMetrics is the span name for getting metrics.
	SpanGetMetrics = "admin.get_metrics"

	// SpanGetInfo is the span name for getting application info.
	SpanGetInfo = "admin.get_info"

	// SpanUpdateUserRoles is the span name for updating user roles.
	SpanUpdateUserRoles = "admin.update_user_roles"

	// SpanPromoteToAdmin is the span name for promoting a user to admin.
	SpanPromoteToAdmin = "admin.promote_to_admin"

	// SpanDemoteFromAdmin is the span name for demoting a user from admin.
	SpanDemoteFromAdmin = "admin.demote_from_admin"

	// SpanBlockUser is the span name for blocking a user.
	SpanBlockUser = "admin.block_user"

	// SpanUnblockUser is the span name for unblocking a user.
	SpanUnblockUser = "admin.unblock_user"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// SpanEventValidateRoles is the event name for validating roles.
	SpanEventValidateRoles = "ValidateRoles"

	// SpanEventGetUser is the event name for getting user.
	SpanEventGetUser = "GetUser"

	// SpanEventUpdateRoles is the event name for updating roles.
	SpanEventUpdateRoles = "UpdateRoles"

	// SpanEventInvalidateRoleCache is the event name for invalidating role cache.
	SpanEventInvalidateRoleCache = "InvalidateRoleCache"

	// SpanEventRevokeSessions is the event name for revoking user sessions.
	SpanEventRevokeSessions = "RevokeUserSessions"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusHealthy indicates the system is healthy.
	StatusHealthy = "healthy"

	// StatusUnhealthy indicates the system is unhealthy.
	StatusUnhealthy = "unhealthy"

	// StatusRolesUpdated indicates the roles were updated successfully.
	StatusRolesUpdated = "roles_updated"

	// StatusUnauthorized indicates an unauthorized operation.
	StatusUnauthorized = "unauthorized"

	// StatusUserNotFound indicates a user was not found.
	StatusUserNotFound = "user_not_found"

	// StatusUpdateFailed indicates an update operation failed.
	StatusUpdateFailed = "update_failed"

	// StatusPromotedToAdmin indicates a user was promoted to admin.
	StatusPromotedToAdmin = "promoted_to_admin"

	// StatusDemotedFromAdmin indicates a user was demoted from admin.
	StatusDemotedFromAdmin = "demoted_from_admin"

	// StatusUserBlocked indicates a user was blocked.
	StatusUserBlocked = "user_blocked"

	// StatusUserUnblocked indicates a user was unblocked.
	StatusUserUnblocked = "user_unblocked"
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrorNoRolesToUpdate indicates there were no roles to update.
	ErrorNoRolesToUpdate = "no roles to update"

	// ErrorInvalidRole indicates an invalid role was provided.
	ErrorInvalidRole = "invalid role provided"

	// ErrorUserNotFound indicates the user was not found.
	ErrorUserNotFound = "user not found"

	// ErrorToUpdateUserRoles indicates an error when updating user roles.
	ErrorToUpdateUserRoles = "error to update user roles"

	// ErrorToGetUser indicates an error when getting user.
	ErrorToGetUser = "error to get user"

	// ErrorCannotBlockAdmin indicates attempt to block an admin user.
	ErrorCannotBlockAdmin = "cannot block an admin user"

	// ErrorCannotRemoveLastAdmin indicates attempt to remove the last admin.
	ErrorCannotRemoveLastAdmin = "cannot remove admin role from the last admin"

	// ErrorUnauthorizedPromoteToAdmin indicates insufficient privileges to promote.
	ErrorUnauthorizedPromoteToAdmin = "insufficient privileges: only owner or admin can promote to admin"

	// ErrorUnauthorizedDemoteFromAdmin indicates insufficient privileges to demote.
	ErrorUnauthorizedDemoteFromAdmin = "insufficient privileges: only owner or admin can demote from admin"

	// ErrorUnauthorizedBlockUser indicates insufficient privileges to block a user.
	ErrorUnauthorizedBlockUser = "insufficient privileges: cannot block user with equal or higher privilege"

	// ErrorUnauthorizedUnblockUser indicates insufficient privileges to unblock a user.
	ErrorUnauthorizedUnblockUser = "insufficient privileges to unblock user"

	// ErrorPromoteToAdminFailed indicates an error when promoting a user.
	ErrorPromoteToAdminFailed = "failed to promote user"

	// ErrorDemoteFromAdminFailed indicates an error when demoting a user.
	ErrorDemoteFromAdminFailed = "failed to demote user"

	// ErrorBlockUserFailed indicates an error when blocking a user.
	ErrorBlockUserFailed = "failed to block user"

	// ErrorUnblockUserFailed indicates an error when unblocking a user.
	ErrorUnblockUserFailed = "failed to unblock user"
)

// =============================================================================
// BUSINESS LOGIC - Success Messages
// =============================================================================

const (
	// SuccessRolesUpdated indicates the roles were updated successfully.
	SuccessRolesUpdated = "user roles updated successfully"
)

// =============================================================================
// LOGGING - Info Messages
// =============================================================================

const (
	// InfoUpdatingUserRoles is an info message when updating user roles.
	InfoUpdatingUserRoles = "updating user roles"

	// InfoUserRolesUpdated is an info message when user roles were updated.
	InfoUserRolesUpdated = "user roles updated"

	// InfoUserPromotedToAdmin is an info message when a user is promoted to admin.
	InfoUserPromotedToAdmin = "user promoted to admin"

	// InfoUserDemotedFromAdmin is an info message when a user is demoted from admin.
	InfoUserDemotedFromAdmin = "user demoted from admin"

	// InfoUserBlocked is an info message when a user is blocked.
	InfoUserBlocked = "user blocked"

	// InfoUserUnblocked is an info message when a user is unblocked.
	InfoUserUnblocked = "user unblocked"
)

// =============================================================================
// LOGGING - Warning Messages
// =============================================================================

const (
	// WarnInvalidRoleProvided is a warning when an invalid role is provided.
	WarnInvalidRoleProvided = "invalid role provided in request"

	// WarnUnauthorizedPromoteAttempt is a warning when a promote attempt is unauthorized.
	WarnUnauthorizedPromoteAttempt = "unauthorized promote attempt"

	// WarnUnauthorizedBlockAttempt is a warning when a block attempt is unauthorized.
	WarnUnauthorizedBlockAttempt = "unauthorized block attempt"

	// WarnFailedToInvalidateRoleCache is a warning when role cache invalidation fails.
	WarnFailedToInvalidateRoleCache = "failed to invalidate role cache"

	// WarnFailedToRevokeUserSessions is a warning when session revocation fails.
	WarnFailedToRevokeUserSessions = "failed to revoke user sessions"
)

// =============================================================================
// LOGGING - Keys
// =============================================================================

const (
	// LogKeyActorUserID is the log key for the actor user ID.
	LogKeyActorUserID = "actor_user_id"

	// LogKeyByUserID is the log key for the actor user ID used in audit logs.
	LogKeyByUserID = "by_user_id"

	// LogKeyActorRole is the log key for the actor role.
	LogKeyActorRole = "actor_role"

	// LogKeyTargetUserID is the log key for the target user ID.
	LogKeyTargetUserID = "target_user_id"

	// LogKeyTargetRole is the log key for the target user role.
	LogKeyTargetRole = "target_role"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrNoRolesToUpdate is a sentinel error when no roles are provided for update.
	ErrNoRolesToUpdate = errors.New(ErrorNoRolesToUpdate)

	// ErrInvalidRole is a sentinel error for invalid role.
	ErrInvalidRole = errors.New(ErrorInvalidRole)

	// ErrUserNotFound is a sentinel error when user is not found.
	ErrUserNotFound = errors.New(ErrorUserNotFound)

	// ErrUpdateUserRoles is a sentinel error for role update failures.
	ErrUpdateUserRoles = errors.New(ErrorToUpdateUserRoles)

	// ErrGetUser is a sentinel error for getting user failures.
	ErrGetUser = errors.New(ErrorToGetUser)

	// ErrCannotBlockAdmin is a sentinel error when trying to block an admin.
	ErrCannotBlockAdmin = errors.New(ErrorCannotBlockAdmin)

	// ErrCannotRemoveLastAdmin is a sentinel error when removing the last admin.
	ErrCannotRemoveLastAdmin = errors.New(ErrorCannotRemoveLastAdmin)

	// ErrUnauthorizedPromoteToAdmin is a sentinel error for unauthorized promotion.
	ErrUnauthorizedPromoteToAdmin = errors.New(ErrorUnauthorizedPromoteToAdmin)

	// ErrUnauthorizedDemoteFromAdmin is a sentinel error for unauthorized demotion.
	ErrUnauthorizedDemoteFromAdmin = errors.New(ErrorUnauthorizedDemoteFromAdmin)

	// ErrUnauthorizedBlockUser is a sentinel error for unauthorized block.
	ErrUnauthorizedBlockUser = errors.New(ErrorUnauthorizedBlockUser)

	// ErrUnauthorizedUnblockUser is a sentinel error for unauthorized unblock.
	ErrUnauthorizedUnblockUser = errors.New(ErrorUnauthorizedUnblockUser)

	// ErrPromoteToAdminFailed is a sentinel error for promote failures.
	ErrPromoteToAdminFailed = errors.New(ErrorPromoteToAdminFailed)

	// ErrDemoteFromAdminFailed is a sentinel error for demote failures.
	ErrDemoteFromAdminFailed = errors.New(ErrorDemoteFromAdminFailed)

	// ErrBlockUserFailed is a sentinel error for block failures.
	ErrBlockUserFailed = errors.New(ErrorBlockUserFailed)

	// ErrUnblockUserFailed is a sentinel error for unblock failures.
	ErrUnblockUserFailed = errors.New(ErrorUnblockUserFailed)
)
