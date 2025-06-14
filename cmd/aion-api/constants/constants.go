// Package constants contains constants used throughout the application.
package constants

// StartingApplication is a constant string used to indicate the application is starting its initialization process.
const StartingApplication = "starting application"

// SuccessToLoadConfiguration is a constant string indicating that the configuration has been successfully loaded.
const SuccessToLoadConfiguration = "configuration loaded successfully"

// SuccessToInitializeDependencies is a constant string indicating successful initialization of application dependencies.
const SuccessToInitializeDependencies = "dependencies initialized successfully"

// ErrToFailedLoadConfiguration is a constant string representing an error message for failure in loading configuration.
const ErrToFailedLoadConfiguration = "failed to load configuration"

// ErrInitializeDependencies is a constant string representing an error message when dependencies fail to initialize.
const ErrInitializeDependencies = "failed to initialize dependencies"

// ErrStartHTTPServer is a constant string representing an error indicating the failure to start the HTTP server.
const ErrStartHTTPServer = "failed to start server"

// ErrStartGraphqlServer is a constant string used to denote a failure in starting the GraphQL server.
const ErrStartGraphqlServer = "failed to start graphql server"

// ServerHTTPStarted is a constant string indicating that the HTTP server has started successfully.
const ServerHTTPStarted = "server http started"

// GraphqlServerStarted is a constant string that indicates the GraphQL server has started.
const GraphqlServerStarted = "graphql server started"

// MsgShutdownSignalReceived is a constant string logged when the application receives a shutdown signal to start graceful shutdown procedures.
const MsgShutdownSignalReceived = "shutdown signal received, attempting graceful shutdown"

// ErrHTTPGracefulShutdown is a constant string representing an error that occurs during the HTTP server's graceful shutdown process.
const ErrHTTPGracefulShutdown = "error during graceful shutdown"

// ErrGraphqlGracefulShutdown is a constant string representing an error message during the graceful shutdown of the GraphQL server.
const ErrGraphqlGracefulShutdown = "error during graphql server graceful shutdown"

// MsgGracefulShutdownSuccess indicates that the server has been successfully shut down in a graceful manner.
const MsgGracefulShutdownSuccess = "server shutdown gracefully"

// Port represents the key used to define or identify a port configuration parameter.
const Port = "port"

// ContextPath is a constant string representing the base path used for contextual application configuration or routing.
const ContextPath = "contextPath"

// Error is a constant string representing a generic error identifier or key for logging and error handling purposes.
const Error = "error"
