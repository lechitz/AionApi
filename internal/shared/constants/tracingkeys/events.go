// Package tracingkeys defines shared keys for OpenTelemetry span attributes and audit fields.
package tracingkeys

// =============================================================================
// COMMON EVENTSy
// Format: <domain>.<action>.<detail>
// =============================================================================

// Common operation events (reusable across domains).
const (
	EventInputValidate    = "input.validate"
	EventUniquenessCheck  = "uniqueness.check"
	EventRepositoryCreate = "repository.create"
	EventRepositoryGet    = "repository.get"
	EventRepositoryUpdate = "repository.update"
	EventRepositoryDelete = "repository.delete"
	EventRepositoryList   = "repository.list"
	EventCacheGet         = "cache.get"
	EventCacheSet         = "cache.set"
	EventCacheDelete      = "cache.delete"
	EventSuccess          = "success"
	EventError            = "error"
)

// Auth events.
const (
	EventAuthUserLookup      = "auth.user.lookup"
	EventAuthPasswordCompare = "auth.password.compare"
	EventAuthTokenGenerate   = "auth.token.generate" //nolint:gosec // event name, not a credential
	EventAuthTokenValidate   = "auth.token.validate" //nolint:gosec // event name, not a credential
	EventAuthTokenRevoke     = "auth.token.revoke"   //nolint:gosec // event name, not a credential
	EventAuthTokenSave       = "auth.token.save"     //nolint:gosec // event name, not a credential
	EventAuthTokenGet        = "auth.token.get"      //nolint:gosec // event name, not a credential
	EventAuthTokenCompare    = "auth.token.compare"  //nolint:gosec // event name, not a credential
	EventAuthExtractUserID   = "auth.user_id.extract"
	EventAuthLoginSuccess    = "auth.login.success"
)
