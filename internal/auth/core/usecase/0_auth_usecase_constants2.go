// Package constants contains constants related to token operations.
package usecase

// SpanValidateToken is the name of the span for validating a token.
const SpanValidateToken = "ValidateToken"

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

// Events (trace)
const (
	// EventVerifyToken is emitted right before verifying token signature/exp.
	EventVerifyToken = "verify_token"

	// EventExtractUserID is emitted before extracting/parsing the userID claim.
	EventExtractUserID = "extract_user_id"

	// EventGetTokenFromStore is emitted before fetching the token from the cache.
	EventGetTokenFromStore = "get_token_from_store"

	// EventCompareToken is emitted before comparing the provided token with the stored one.
	EventCompareToken = "compare_token"

	// EventTokenValidated is emitted after successful validation.
	EventTokenValidated = "token_validated"
)
