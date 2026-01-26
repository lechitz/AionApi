// Package cache provides role cache operations.
package cache

import "time"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// RoleCacheTracerName is the name of the tracer for role cache operations.
const RoleCacheTracerName = "aionapi.auth.role_cache"

// Span names.
const (
	SpanNameRoleSave   = "auth.role_cache.save"
	SpanNameRoleGet    = "auth.role_cache.get"
	SpanNameRoleDelete = "auth.role_cache.delete"
)

// Operations.
const (
	OperationRoleSave   = "save"
	OperationRoleGet    = "get"
	OperationRoleDelete = "delete"
)

// Attribute keys for tracing and logging.
const (
	AttributeRoleCacheKey = "role_cache_key"
	AttributeRoleTTL      = "role_cache_ttl"
)

// =============================================================================
// BUSINESS LOGIC - Configuration and Messages
// =============================================================================

// RoleExpirationDefault defines the default expiration for role cache (24 hours).
const RoleExpirationDefault = 24 * time.Hour

// RoleKeyFormat : auth:roles:{userID}.
const RoleKeyFormat = "auth:roles:%d"

// Success messages.
const (
	RolesRetrievedSuccessfully = "roles retrieved successfully from cache"
	RolesSavedSuccessfully     = "roles saved successfully to cache"
	RolesDeletedSuccessfully   = "roles deleted successfully from cache"
)

// Error messages.
const (
	ErrorToSaveRolesToCache     = "error to save roles to cache"
	ErrorToGetRolesFromCache    = "error to get roles from cache"
	ErrorToDeleteRolesFromCache = "error to delete roles from cache"
	ErrorToSerializeRoles       = "error to serialize roles"
	ErrorToDeserializeRoles     = "error to deserialize roles"
)
