package main

const (
	// MsgConfigLoaded is the message for when the configuration is loaded.
	MsgConfigLoaded = "configuration loaded"

	// MsgDepsInitialized is the message for when dependencies are initialized.
	MsgDepsInitialized = "dependencies initialized"

	// ErrLoadConfig is the error message for when the configuration fails to load.
	ErrLoadConfig = "failed to load configuration"

	// ErrInvalidConfig is the error message for when the configuration is invalid.
	ErrInvalidConfig = "invalid configuration"

	// ErrInitDeps is the error message for when dependencies fail to initialize.
	ErrInitDeps = "failed to initialize dependencies"

	// ErrServerRunFailed is the error message for when the server fails to run.
	ErrServerRunFailed = "server run failed"
)
