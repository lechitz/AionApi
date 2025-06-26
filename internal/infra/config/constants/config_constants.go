// Package constants contains constants used throughout the application.
package constants

// ErrFailedToProcessEnvVars is returned when environment variables cannot be processed.
const ErrFailedToProcessEnvVars = "failed to process environment variables: %v"

// ErrGenerateSecretKey is returned when the secret key generation fails.
const ErrGenerateSecretKey = "failed to generate secret key"

// SecretKeyWasNotSet is logged when no SECRET_KEY is found, and a new one is generated.
const SecretKeyWasNotSet = "SECRET_KEY was not set. A new one was generated for this runtime session." // #nosec G101
