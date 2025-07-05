// Package constants contains constants used throughout the application.
package constants

// SuccessToLoadConfiguration is a constant string indicating that the configuration has been successfully loaded.
const SuccessToLoadConfiguration = "configuration loaded successfully"

// SuccessToInitializeDependencies is a constant string indicating successful initialization of application dependencies.
const SuccessToInitializeDependencies = "dependencies initialized successfully"

// LoadedConfig is a constant string used to indicate the application has successfully loaded its configuration.
const LoadedConfig = "loaded config: %+v"

// ErrToFailedLoadConfiguration is a constant string representing an error message for failure in loading configuration.
const ErrToFailedLoadConfiguration = "failed to load configuration"

// ErrFailedToStartHTTPServer is a constant string used to indicate a failure to start the HTTP server.
const ErrFailedToStartHTTPServer = "failed to start HTTP server: %w"

// ErrFailedToStartGraphqlServer is a constant string used to indicate a failure to start the GraphQL server.
const ErrFailedToStartGraphqlServer = "failed to start GraphQL server: %w"

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

// MsgUnexpectedServerFailure is a constant string used to indicate an unexpected failure in starting one of the application servers.
const MsgUnexpectedServerFailure = "unexpected failure while starting one of the application servers (HTTP or GraphQL)"
