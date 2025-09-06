// Package constants contains constants used throughout the application.
package main

const (

	// SuccessToLoadConfiguration is a constant string indicating that the configuration has been successfully loaded.
	SuccessToLoadConfiguration = "configuration loaded successfully"

	// SuccessToInitializeDependencies is a constant string indicating successful initialization of application dependencies.
	SuccessToInitializeDependencies = "dependencies initialized successfully"

	// ErrToFailedLoadConfiguration is a constant string representing an error message for failure in loading configuration.
	ErrToFailedLoadConfiguration = "failed to load configuration"

	// ErrInvalidConfiguration is a constant string used to indicate an invalid configuration.
	ErrInvalidConfiguration = "invalid configuration"

	// ErrInitializeDependencies is a constant string representing an error message when dependencies fail to initialize.
	ErrInitializeDependencies = "failed to initialize dependencies"

	// ErrStartHTTPServer is a constant string representing an error indicating the failure to start the HTTP server.
	ErrStartHTTPServer = "failed to start server" //

	// ErrStartGraphqlServer is a constant string used to denote a failure in starting the GraphQL server.
	ErrStartGraphqlServer = "failed to start graphql server"

	// MsgShutdownSignalReceived is a constant string logged when the application receives a shutdown signal to start graceful shutdown procedures.
	MsgShutdownSignalReceived = "shutdown signal received, attempting graceful shutdown"

	// MsgUnexpectedServerFailure is a constant string used to indicate an unexpected failure in starting one of the application servers.
	MsgUnexpectedServerFailure = "unexpected failure while starting one of the application servers (HTTP or GraphQL)"

	// ServerFailureFmt is a constant string used to format the message indicating a failure in starting the server.
	ServerFailureFmt = "failed to start server on %s: %w"

	// ShutdownFailureFmt is a constant string used to format the message indicating a failure in shutting down the server.
	ShutdownFailureFmt = "failed to shutdown server"
)
