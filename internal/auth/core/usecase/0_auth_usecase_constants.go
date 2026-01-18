package usecase

import (
	"errors"
	"time"
)

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

	// SpanGenerateAndStoreTokens is the span name for generating and storing tokens.
	SpanGenerateAndStoreTokens = "auth.token.generate_and_store" // #nosec G101: span name, not a credential

	// SpanRefreshTokenRenewal is the span name for refresh token renewal.
	SpanRefreshTokenRenewal = "auth.token.refresh" // #nosec G101: span name, not a credential
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

	// EventRevokeAccessToken is emitted before revoking the token.
	EventRevokeAccessToken = "auth.token.revoke.access_token" // #nosec G101: event name, not a credential

	// EventRevokeRefreshToken is emitted before revoking the refresh token.
	EventRevokeRefreshToken = "auth.token.revoke.refresh_token" // #nosec G101: event name, not a credential

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

	// EventCacheUserProfile is emitted before caching user profile on login.
	EventCacheUserProfile = "auth.user.cache_profile"

	// EventGenerateRefreshToken is emitted before generating refresh token.
	EventGenerateRefreshToken = "auth.token.generate.refresh" // #nosec G101: event name, not a credential

	// EventCheckGracePeriod is emitted when primary token doesn't match and grace period will be checked.
	EventCheckGracePeriod = "auth.token.check_grace_period" // #nosec G101: event name, not a credential

	// EventValidatedViaGrace is emitted when token validation succeeds via grace period.
	EventValidatedViaGrace = "auth.token.validated_via_grace" // #nosec G101: event name, not a credential

	// EventValidateRefreshToken is emitted before validating the refresh token and checking the store.
	EventValidateRefreshToken = "auth.token.refresh.validate" // #nosec G101: event name, not a credential

	// EventGenerateAndSaveNewTokens is emitted before generating and persisting new access/refresh tokens.
	EventGenerateAndSaveNewTokens = "auth.token.refresh.generate_and_save" // #nosec G101: event name, not a credential

	// EventSaveAccessTokenToStore is emitted before saving the access token to the store.
	EventSaveAccessTokenToStore = "auth.token.save.access" // #nosec G101: event name, not a credential

	// EventSaveRefreshTokenToStore is emitted before saving the refresh token to the store.
	EventSaveRefreshTokenToStore = "auth.token.save.refresh" // #nosec G101: event name, not a credential
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

	// ErrorInvalidRefreshToken is the error message when refresh token verification or store match fails.
	ErrorInvalidRefreshToken = "invalid refresh token"
)

// Success messages.
const (
	// SuccessToLogin is the success message when the user logs in.
	SuccessToLogin = "user logged in successfully"

	// SuccessUserLoggedOut is the success message when the user logs out.
	SuccessUserLoggedOut = "user logged out successfully"
)

// Log messages.
const (
	// LogFailedToCacheUserProfile is logged when user profile caching fails on login.
	LogFailedToCacheUserProfile = "failed to cache user profile on login"

	// LogFailedToCacheUserData is logged when user caching fails during refresh token renewal.
	LogFailedToCacheUserData = "failed to cache user data"

	// =============================================================================
	// Refresh/Grace period logs.
	// =============================================================================.

	// LogNoPreviousTokenForGrace is logged when no previous token is found for grace period.
	LogNoPreviousTokenForGrace = "no previous token found for grace period"

	// LogSkippingGraceNoOldToken is logged when grace period is skipped due to no old token.
	LogSkippingGraceNoOldToken = "skipping grace period: no old token"

	// LogSkippingGraceTokensIdentical is logged when grace period is skipped due to identical tokens.
	LogSkippingGraceTokensIdentical = "skipping grace period: tokens are identical"

	// LogMovingTokenToGrace is logged when moving old token to grace period.
	LogMovingTokenToGrace = "moving token to grace period"

	// LogFailedToSaveTokenToGrace is logged when saving token to grace period fails.
	LogFailedToSaveTokenToGrace = "failed to save token to grace period"

	// LogTokenMovedToGraceSuccess is logged when token is moved to grace period successfully.
	LogTokenMovedToGraceSuccess = "token moved to grace period successfully"

	// =============================================================================
	// Validate/grace period logs.
	// =============================================================================.

	LogTokenMismatchCheckingGrace = "token mismatch with primary, checking grace period"
	LogTokenValidatedViaGrace     = "token validated via grace period"
	LogTokenNotFoundInGrace       = "token not found in grace period"

	// =============================================================================
	// Log/trace operation name for refresh token renewal.
	// =============================================================================.

	// OperationRefreshTokenRenewal is the operation name for refresh token renewal.
	OperationRefreshTokenRenewal = "RefreshTokenRenewal"
)

// =============================================================================
// BUSINESS LOGIC - Refresh Token
// =============================================================================

const (
	// SuccessRefreshTokenRenewed indicates refresh token renewal succeeded.
	SuccessRefreshTokenRenewed = "refresh token renewed successfully"

	// AuthGraceKeyPrefix Redis key prefix for grace period tokens.
	AuthGraceKeyPrefix = "auth:grace"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrInvalidToken is a sentinel error for invalid tokens.
	ErrInvalidToken = errors.New(ErrorInvalidToken)

	// ErrInvalidUserIDClaim is a sentinel error for invalid user ID claims.
	ErrInvalidUserIDClaim = errors.New(ErrorInvalidUserIDClaim)

	// ErrTokenRetrievalFromCache is a sentinel error for cache retrieval failures.
	ErrTokenRetrievalFromCache = errors.New(ErrorToRetrieveTokenFromCache)

	// ErrTokenMismatch is a sentinel error for token mismatch.
	ErrTokenMismatch = errors.New(ErrorTokenMismatch)

	// ErrInvalidCredentials is a sentinel error for invalid credentials.
	ErrInvalidCredentials = errors.New(InvalidCredentials)

	// ErrUserNotFoundOrInvalidCredentials is a sentinel error for user not found or invalid credentials.
	ErrUserNotFoundOrInvalidCredentials = errors.New(UserNotFoundOrInvalidCredentials)

	// ErrTokenCreation is a sentinel error for token creation failures.
	ErrTokenCreation = errors.New(ErrorToCreateToken)

	// ErrTokenDeletion is a sentinel error for token deletion failures.
	ErrTokenDeletion = errors.New(ErrorToDeleteToken)

	// ErrGetUserByUsername is a sentinel error for user retrieval failures.
	ErrGetUserByUsername = errors.New(ErrorToGetUserByUserName)

	// ErrInvalidRefreshToken is a sentinel error for invalid refresh token.
	ErrInvalidRefreshToken = errors.New(ErrorInvalidRefreshToken)
)
