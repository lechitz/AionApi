package constants

const (
	Settings                  = "setting"
	ErrFailedToProcessEnvVars = "failed to process environment variables: %v"
	ErrServerContextEmpty     = "server context path SERVER_CONTEXT must be configured and cannot be empty"
	ErrServerPortEmpty        = "server port SERVER_PORT must be configured and cannot be empty"
	ErrGenerateSecretKey      = "failed to generate secret key"
	SecretKeyWasNotSet        = "SECRET_KEY was not set. A new one was generated for this runtime session."
)
