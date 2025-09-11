// cmd/aion-api/0_cmd_aion-api_constants.go
package main

const (
	// Mensagens de sucesso (logs informativos)
	MsgConfigLoaded    = "configuration loaded"
	MsgDepsInitialized = "dependencies initialized"

	// Mensagens de erro (logs + os.Exit(1) na main)
	ErrLoadConfig      = "failed to load configuration"
	ErrInvalidConfig   = "invalid configuration"
	ErrInitDeps        = "failed to initialize dependencies"
	ErrServerRunFailed = "server run failed"
)
