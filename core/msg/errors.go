package msg

const StartingApplication = "starting application"

const (
	// Errors related to server startup

	ErrInitializeDependencies    = "failed to initialize application dependencies"
	ErrStartServer               = "failed to start HTTP server"
	ErrToFailedLoadConfiguration = "failed to load configuration: %v"
	ErrToStartServer             = "failed to start server: %v"

	ErrInvalidConfiguration         = "invalid configuration"
	ErrorMissingServerConfiguration = "missing server configuration"
	ErrServerContextEmpty           = "server context path SERVER_CONTEXT must be configured and cannot be empty"
	ErrFailedToProcessEnvVars       = "failed to process environment variables: %v"
	ErrFailedToLoadEnvFile          = "failed to load environment file: %v"
	ErrServerPortEmpty              = "server port SERVER_PORT must be configured and cannot be empty"
	ErrSecretKeyEmpty               = "secret key SECRET_KEY must be configured and cannot be empty"
	ErrGenerateSecretKey            = "failed to generate secret key"
)

// Success messages for server startup

const (
	ServerStarted                   = "server started"
	SuccessToLoadConfiguration      = "successfully loaded configuration"
	SuccessToInitializeDependencies = "successfully initialized dependencies"
)
