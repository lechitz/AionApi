package constants

const (
	StartingApplication             = "starting application"
	ErrToFailedLoadConfiguration    = "failed to load configuration: %v"
	ErrInitializeDependencies       = "failed to initialize application dependencies"
	SuccessToInitializeDependencies = "successfully initialized dependencies"
	ErrStartServer                  = "failed to start HTTP server"
	ServerStarted                   = "server started"
	Error                           = "error"
	ContextPath                     = "contextPath"
	Port                            = "port"
	SuccessToLoadConfiguration      = "successfully loaded configuration"
)

const (
	Settings                  = "setting"
	ErrFailedToProcessEnvVars = "failed to process environment variables: %v"
	ErrServerContextEmpty     = "server context path SERVER_CONTEXT must be configured and cannot be empty"
	ErrServerPortEmpty        = "server port SERVER_PORT must be configured and cannot be empty"
	ErrGenerateSecretKey      = "failed to generate secret key"
	SecretKeyWasNotSet        = "SECRET_KEY was not set. A new one was generated for this runtime session."
)
