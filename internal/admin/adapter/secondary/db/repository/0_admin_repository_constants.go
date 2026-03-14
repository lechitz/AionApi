// Package repository holds repository-scoped constants to avoid magic strings.
package repository

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer used in the admin repository.
// Format: aionapi.<domain>.<layer>.
const TracerName = "aionapi.admin.repository"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanGetByID is the span name for getting by ID.
	SpanGetByID = "admin.repository.get_by_id"

	// SpanUpdateRoles is the span name for updating roles.
	SpanUpdateRoles = "admin.repository.update_roles"

	// SpanAssignDefaultRole is the span name for assigning default role.
	SpanAssignDefaultRole = "admin.repository.assign_default_role"
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// OperationGetByID is the attribute value for getting a user by ID.
	OperationGetByID = "get_by_id"

	// OperationUpdateRoles is the attribute value for updating user roles.
	OperationUpdateRoles = "update_roles"

	// OperationAssignDefaultRole is the attribute value for assigning default role.
	OperationAssignDefaultRole = "assign_default_role"
)

// =============================================================================
// BUSINESS LOGIC - Status
// =============================================================================

const (
	// StatusUserFetched indicates a user was fetched successfully.
	StatusUserFetched = "user_fetched"

	// StatusRolesUpdated indicates roles were updated successfully.
	StatusRolesUpdated = "roles_updated"

	// StatusDefaultRoleAssigned indicates default role was assigned successfully.
	StatusDefaultRoleAssigned = "default_role_assigned"
)

// =============================================================================
// LOGGING - Messages
// =============================================================================

const (
	// LogFetchingUser is the log message for fetching a user.
	LogFetchingUser = "fetching user by ID"

	// LogUserFetched is the log message when user is fetched.
	LogUserFetched = "user fetched successfully"

	// LogUpdatingRoles is the log message for updating roles.
	LogUpdatingRoles = "updating user roles"

	// LogRolesUpdated is the log message when roles are updated.
	LogRolesUpdated = "user roles updated successfully"

	// LogFailedGetUser is the log message for failed user fetch.
	LogFailedGetUser = "failed to get user"

	// LogFailedUpdateRoles is the log message for failed roles update.
	LogFailedUpdateRoles = "failed to update user roles"

	// LogAssigningDefaultRole is the log message for assigning default role.
	LogAssigningDefaultRole = "assigning default role to user"

	// LogDefaultRoleAssigned is the log message when default role is assigned.
	LogDefaultRoleAssigned = "default role assigned successfully"

	// LogFailedToAssignDefaultRole is the log message for failed default role assignment.
	LogFailedToAssignDefaultRole = "failed to assign default role"
)

// =============================================================================
// DATABASE - Queries and Values
// =============================================================================

const (
	// DefaultRoleName is the name of the default role assigned to new users.
	DefaultRoleName = "user"

	// QueryCountDefaultRole counts active default roles.
	QueryCountDefaultRole = `SELECT COUNT(*) FROM aion_api.roles WHERE name = $1 AND is_active = true`

	// QueryInsertDefaultRoleAssignment inserts a default role assignment.
	QueryInsertDefaultRoleAssignment = `
		INSERT INTO aion_api.user_roles (user_id, role_id, assigned_at)
		SELECT $1, role_id, $2
		FROM aion_api.roles
		WHERE name = $3 AND is_active = true
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

const (
	// ErrorCheckDefaultRoleExists indicates an error checking the default role existence.
	ErrorCheckDefaultRoleExists = "failed to check if 'user' role exists"

	// ErrorDefaultRoleNotFound indicates the default role was not found.
	ErrorDefaultRoleNotFound = "default role 'user' not found in database"

	// ErrorAssignDefaultRole indicates an error assigning the default role.
	ErrorAssignDefaultRole = "failed to assign default role"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrCheckDefaultRoleExists is a sentinel error for default role existence checks.
	ErrCheckDefaultRoleExists = errors.New(ErrorCheckDefaultRoleExists)

	// ErrDefaultRoleNotFound is a sentinel error when the default role is missing.
	ErrDefaultRoleNotFound = errors.New(ErrorDefaultRoleNotFound)

	// ErrAssignDefaultRole is a sentinel error for default role assignment failures.
	ErrAssignDefaultRole = errors.New(ErrorAssignDefaultRole)
)
