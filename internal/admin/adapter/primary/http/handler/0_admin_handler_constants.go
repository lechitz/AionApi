// Package handler is the handler for the admin context in the application.
package handler

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the tracer name for admin handler.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.admin.handler"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanHealthCheck is the span name for health check handler.
	SpanHealthCheck = "admin.handler.health_check"

	// SpanGetMetrics is the span name for metrics handler.
	SpanGetMetrics = "admin.handler.get_metrics"

	// SpanGetInfo is the span name for info handler.
	SpanGetInfo = "admin.handler.get_info"

	// SpanUpdateUserRoles is the span name for update user roles handler.
	SpanUpdateUserRoles = "admin.handler.update_user_roles"

	// SpanPromoteToAdmin is the span name for promote to admin handler.
	SpanPromoteToAdmin = "admin.handler.promote_to_admin"

	// SpanDemoteFromAdmin is the span name for demote from admin handler.
	SpanDemoteFromAdmin = "admin.handler.demote_from_admin"

	// SpanBlockUser is the span name for block user handler.
	SpanBlockUser = "admin.handler.block_user"

	// SpanUnblockUser is the span name for unblock user handler.
	SpanUnblockUser = "admin.handler.unblock_user"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// EventDecodeRequest is the event name for decoding request.
	EventDecodeRequest = "decode_request"

	// EventUserIDExtracted is the event name when user ID is extracted from path.
	EventUserIDExtracted = "user_id_extracted"

	// EventAdminServiceUpdateRoles is the event name for calling admin service.
	EventAdminServiceUpdateRoles = "admin_service.update_user_roles"

	// EventRolesUpdatedSuccess is the event name for successful roles update.
	EventRolesUpdatedSuccess = "roles.updated.success"

	// EventAdminServicePromoteToAdmin is the event name for calling promote to admin service.
	EventAdminServicePromoteToAdmin = "admin_service.promote_to_admin"

	// EventPromoteToAdminSuccess is the event name for successful promotion to admin.
	EventPromoteToAdminSuccess = "roles.promote_to_admin.success"

	// EventAdminServiceDemoteFromAdmin is the event name for calling demote from admin service.
	EventAdminServiceDemoteFromAdmin = "admin_service.demote_from_admin"

	// EventDemoteFromAdminSuccess is the event name for successful demotion from admin.
	EventDemoteFromAdminSuccess = "roles.demote_from_admin.success"

	// EventAdminServiceBlockUser is the event name for calling block user service.
	EventAdminServiceBlockUser = "admin_service.block_user"

	// EventBlockUserSuccess is the event name for successful user block.
	EventBlockUserSuccess = "roles.block_user.success"

	// EventAdminServiceUnblockUser is the event name for calling unblock user service.
	EventAdminServiceUnblockUser = "admin_service.unblock_user"

	// EventUnblockUserSuccess is the event name for successful user unblock.
	EventUnblockUserSuccess = "roles.unblock_user.success"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusHealthy indicates the system is healthy.
	StatusHealthy = "healthy"

	// StatusUnhealthy indicates the system is unhealthy.
	StatusUnhealthy = "unhealthy"

	// StatusRolesUpdated indicates roles were updated successfully.
	StatusRolesUpdated = "roles_updated"

	// StatusPromotedToAdmin indicates user was promoted to admin.
	StatusPromotedToAdmin = "promoted_to_admin"

	// StatusDemotedFromAdmin indicates user was demoted from admin.
	StatusDemotedFromAdmin = "demoted_from_admin"

	// StatusUserBlocked indicates user was blocked.
	StatusUserBlocked = "user_blocked"

	// StatusUserUnblocked indicates user was unblocked.
	StatusUserUnblocked = "user_unblocked"
)

// =============================================================================
// BUSINESS LOGIC - Messages
// =============================================================================

const (
	// MsgHealthCheckSuccess is the message for successful health check.
	MsgHealthCheckSuccess = "health check passed"

	// MsgHealthCheckFailed is the message for failed health check.
	MsgHealthCheckFailed = "health check failed"

	// MsgRolesUpdated is the message for successful roles update.
	MsgRolesUpdated = "user roles updated successfully"

	// MsgUserPromotedToAdmin is the message for successful promotion to admin.
	MsgUserPromotedToAdmin = "user promoted to admin successfully"

	// MsgUserDemotedFromAdmin is the message for successful demotion from admin.
	MsgUserDemotedFromAdmin = "user demoted from admin successfully"

	// MsgUserBlocked is the message for successful user block.
	MsgUserBlocked = "user blocked successfully"

	// MsgUserUnblocked is the message for successful user unblock.
	MsgUserUnblocked = "user unblocked successfully"
)

// =============================================================================
// ERROR MESSAGES
// =============================================================================

const (
	// ErrMissingUserIDParam is the error message for missing user ID parameter.
	ErrMissingUserIDParam = "missing user ID parameter"

	// ErrInvalidUserIDParam is the error message for invalid user ID parameter.
	ErrInvalidUserIDParam = "invalid user ID parameter"

	// ErrUpdateUserRoles is the error message for role update errors.
	ErrUpdateUserRoles = "error updating user roles"

	// ErrMissingActorUserID is the error message for missing actor user ID in context.
	ErrMissingActorUserID = "missing authenticated user ID in context"

	// ErrPromoteToAdminFailed is the error message for promote to admin failures.
	ErrPromoteToAdminFailed = "failed to promote user to admin"

	// ErrDemoteFromAdminFailed is the error message for demote from admin failures.
	ErrDemoteFromAdminFailed = "failed to demote user from admin"

	// ErrBlockUserFailed is the error message for block user failures.
	ErrBlockUserFailed = "failed to block user"

	// ErrUnblockUserFailed is the error message for unblock user failures.
	ErrUnblockUserFailed = "failed to unblock user"

	// LogKeyActorUserID is the log key for actor user ID.
	LogKeyActorUserID = "actor_user_id"
)

// =============================================================================
// SENTINEL ERRORS
// =============================================================================

var (
	// ErrMissingUserID is a sentinel error for missing user ID.
	ErrMissingUserID = errors.New(ErrMissingUserIDParam)

	// ErrInvalidUserID is a sentinel error for invalid user ID.
	ErrInvalidUserID = errors.New(ErrInvalidUserIDParam)
)
