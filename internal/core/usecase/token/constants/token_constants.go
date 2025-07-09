// Package constants contains constants related to token operations.
package constants

// ErrorToSaveToken indicates an error saving the access reference.
const ErrorToSaveToken = "error saving access reference" // #nosec G101

// ErrorToGetToken indicates an error retrieving the access reference.
const ErrorToGetToken = "failed to retrieve token from cache"

// ErrorToDeleteToken indicates an error deleting the access reference.
const ErrorToDeleteToken = "error deleting access reference" // #nosec G101

// ErrorToAssignToken indicates an error assigning the access reference.
const ErrorToAssignToken = "error assigning access reference"

// ErrorInvalidToken indicates the token is invalid.
const ErrorInvalidToken = "invalid access reference"

// ErrorInvalidTokenClaims indicates invalid claims in the token reference.
const ErrorInvalidTokenClaims = "invalid claims in reference"

// ErrorInvalidUserIDClaim indicates the user ID in claims is invalid.
const ErrorInvalidUserIDClaim = "invalid userID in claims"

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
