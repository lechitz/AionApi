package usecase

import "time"

// =============================================================================
// CONFIGURATION
// =============================================================================

// GracePeriodDuration defines how long an old token remains valid after refresh.
// This prevents race conditions in multi-tab scenarios where multiple refreshes
// happen simultaneously. The old token is moved to a grace period cache key
// and remains valid for this duration before expiring.
// Increased to 60s to handle slower tab synchronization and network delays.
const GracePeriodDuration = 60 * time.Second

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName identifies the tracer used in auth use cases.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.auth.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanLogin is the span name for login operation.
	SpanLogin = "auth.login"

	// SpanLogout is the span name for logout operation.
	SpanLogout = "auth.logout"

	// SpanRevokeToken is the span name for token revocation.
	SpanRevokeToken = "auth.token.revoke" // #nosec G101: span name, not a credential

	// SpanValidateToken is the span name for token validation.
	SpanValidateToken = "auth.token.validate" // #nosec G101: span name, not a credential
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// EventLookupUser is emitted before looking up the user.
	EventLookupUser = "auth.user.lookup"

	// EventComparePassword is emitted before comparing the password.
	EventComparePassword = "auth.password.compare" //nolint:goconst // unique event name

	// EventGenerateToken is emitted before generating the token.
	EventGenerateToken = "auth.token.generate" // #nosec G101: event name, not a credential

	// EventSaveTokenToStore is emitted before saving the token to the cache.
	EventSaveTokenToStore = "auth.token.save" // #nosec G101: false positive — event name string, not a credential

	// EventTokenRevoked is emitted after token revocation.
	EventTokenRevoked = "auth.token.revoked" // #nosec G101: event name, not a credential

	// EventLoginSuccess is emitted after successful login.
	EventLoginSuccess = "auth.login.success"

	// EventRevokeToken is emitted before revoking the token.
	EventRevokeToken = "auth.token.revoke" // #nosec G101: event name, not a credential

	// EventVerifyToken is emitted before verifying token signature/exp.
	EventVerifyToken = "auth.token.verify" // #nosec G101: event name, not a credential

	// EventExtractUserID is emitted before extracting/parsing the userID claim.
	EventExtractUserID = "auth.user_id.extract"

	// EventGetTokenFromStore is emitted before fetching the token from the cache.
	EventGetTokenFromStore = "auth.token.get" // #nosec G101: event name, not a credential

	// EventCompareToken is emitted before comparing the provided token with the stored one.
	EventCompareToken = "auth.token.compare" // #nosec G101: event name, not a credential

	// EventTokenValidated is emitted when token validation succeeds.
	EventTokenValidated = "auth.token.validated" // #nosec G101: event name, not a credential
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

// ErrorInvalidToken indicates the token is invalid.
const ErrorInvalidToken = "invalid access reference"

// ErrorInvalidUserIDClaim indicates the user ID in claimsextractor is invalid.
const ErrorInvalidUserIDClaim = "invalid userID in claimsextractor"

// ErrorToRetrieveTokenFromCache indicates an error retrieving the access reference from cache.
const ErrorToRetrieveTokenFromCache = "error retrieving access reference from cache"

// ErrorTokenMismatch indicates the provided token does not match the stored one.
const ErrorTokenMismatch = "provided reference does not match stored one"

// SuccessTokenValidated indicates the access reference was validated successfully.
const SuccessTokenValidated = "access reference validated successfully"

// Errors.
const (
	// ErrorToCompareHashAndPassword is the error message when the password is invalid.
	ErrorToCompareHashAndPassword = "invalid credentials" // #nosec G101

	// ErrorToCreateToken is the error message when the token cannot be created.
	ErrorToCreateToken = "error to create token" // #nosec G101

	// ErrorToDeleteToken is the error message when the token cannot be deleted.
	ErrorToDeleteToken = "error to delete token" // #nosec G101

	// ErrorToGetUserByUserName is the error message when the user cannot be found.
	ErrorToGetUserByUserName = "error to get user by username" // #nosec G101

	// UserNotFoundOrInvalidCredentials is the error message when the user cannot be found.
	UserNotFoundOrInvalidCredentials = "user not found or invalid credentials"

	// InvalidCredentials is the error message when the credentials are invalid.
	InvalidCredentials = "invalid credentials" // #nosec G101
)

// Success messages.
const (
	// SuccessToLogin is the success message when the user logs in.
	SuccessToLogin = "user logged in successfully"

	// SuccessUserLoggedOut is the success message when the user logs out.
	SuccessUserLoggedOut = "user logged out successfully"
)
