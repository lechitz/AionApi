// Package constants contains constants related to token operations.
package constants

// ErrorToSaveToken indicates an error saving the access reference.
const ErrorToSaveToken = "error saving access reference" // #nosec G101

// ErrorToDeleteToken indicates an error deleting the access reference.
const ErrorToDeleteToken = "error deleting access reference" // #nosec G101

// ErrorToAssignToken indicates an error assigning the access reference.
const ErrorToAssignToken = "error assigning access reference"

// ErrorInvalidToken indicates the token is invalid.
const ErrorInvalidToken = "invalid access reference"

// ErrorInvalidTokenClaims indicates invalid claimskeys in the token reference.
const ErrorInvalidTokenClaims = "invalid claimskeys in reference"

// ErrorInvalidUserIDClaim indicates the user ID in claimskeys is invalid.
const ErrorInvalidUserIDClaim = "invalid userID in claimskeys"

// ErrorToRetrieveTokenFromCache indicates an error retrieving the access reference from cache.
const ErrorToRetrieveTokenFromCache = "error retrieving access reference from cache"

// ErrorTokenMismatch indicates the provided token does not match the stored one.
const ErrorTokenMismatch = "provided reference does not match stored one"

// SuccessTokenCreated indicates the access reference was created successfully.
const SuccessTokenCreated = "access reference created successfully"

// SuccessTokenValidated indicates the access reference was validated successfully.
const SuccessTokenValidated = "access reference validated successfully"

// SuccessTokenDeleted indicates the access reference was deleted successfully.
const SuccessTokenDeleted = "access reference deleted successfully"

// TODO: verificar onde as Const abaixo devem ir !

// TokenFromCookie is the key for a token retrieved from cookies.
const TokenFromCookie = "TokenFromCookie"

// TokenFromCache is the key for a token retrieved from cache.
const TokenFromCache = "TokenFromCache"

// SecretKey is the key for secrets used in token operations.
const SecretKey = "secret"

// UserID is the key for user ID values in contexts and payloads.
const UserID = "user_id"
