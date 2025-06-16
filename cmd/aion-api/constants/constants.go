// Package constants contains constants used throughout the application.
package constants

// StartingApplication is a constant string used to indicate the application is starting its initialization process.
const StartingApplication = "starting application"

// SuccessToLoadConfiguration is a constant string indicating that the configuration has been successfully loaded.
const SuccessToLoadConfiguration = "configuration loaded successfully"

// SuccessToInitializeDependencies is a constant string indicating successful initialization of application dependencies.
const SuccessToInitializeDependencies = "dependencies initialized successfully"

// LoadedConfig is a constant string used to indicate the application has successfully loaded its configuration.
const LoadedConfig = "loaded config: %+v"

// ErrToFailedLoadConfiguration is a constant string representing an error message for failure in loading configuration.
const ErrToFailedLoadConfiguration = "failed to load configuration"

// ErrFailedToInitializeOTLPMetricsExporter is a constant string used to indicate a failure to initialize the OTLP metrics exporter.
const ErrFailedToInitializeOTLPMetricsExporter = "failed to initialize OTLP metric exporter"

// ErrInitializeOTPL is a constant string used to indicate a failure to initialize the OTLP exporter.
const ErrInitializeOTPL = "failed to initialize OTLP exporter"

// ErrFailedToShutdownTracerProvider is a constant string used to indicate a failure to shut down the tracer provider.
const ErrFailedToShutdownTracerProvider = "failed to shutdown tracer provider"

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

// Port represents the key used to define or identify a port configuration parameter.
const Port = "port"

// ContextPath is a constant string representing the base path used for contextual application configuration or routing.
const ContextPath = "contextPath"

const GraphQLPath = "/graphql"

// Error is a constant string representing a generic error identifier or key for logging and error handling purposes.
const Error = "error"
