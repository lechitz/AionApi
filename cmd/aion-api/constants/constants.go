package constants

const (
	StartingApplication             = "starting application"
	SuccessToLoadConfiguration      = "configuration loaded successfully"
	SuccessToInitializeDependencies = "dependencies initialized successfully"
	ErrToFailedLoadConfiguration    = "failed to load configuration"
	ErrInitializeDependencies       = "failed to initialize dependencies"
	ErrStartHTTPServer              = "failed to start server"
	ErrStartGraphqlServer           = "failed to start graphql server"
	ServerHTTPStarted               = "server http started"
	GraphqlServerStarted            = "graphql server started"
	MsgShutdownSignalReceived       = "shutdown signal received, attempting graceful shutdown"
	ErrGracefulShutdown             = "error during graceful shutdown"
	MsgGracefulShutdownSuccess      = "server shutdown gracefully"
	Port                            = "port"
	ContextPath                     = "contextPath"
	Error                           = "error"
)
