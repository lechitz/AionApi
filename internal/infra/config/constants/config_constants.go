package constants

const (
	Settings                  = "setting"
	ErrFailedToProcessEnvVars = "failed to process environment variables: %v"
	ErrGenerateSecretKey      = "failed to generate secret key"
	SecretKeyWasNotSet        = "SECRET_KEY was not set. A new one was generated for this runtime session."
)

const (
	SecretKeyFormat = "\nSECRET_KEY=%s\n\n"
)
