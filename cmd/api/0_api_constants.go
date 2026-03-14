package main

import "time"

const (
	defaultStartTimeout = 15 * time.Second
	defaultStopTimeout  = 10 * time.Second

	envBootstrapStartTimeout = "BOOTSTRAP_START_TIMEOUT"
	envBootstrapStopTimeout  = "BOOTSTRAP_STOP_TIMEOUT"

	logMsgBootstrapConfigFailed  = "bootstrap config failed"
	logMsgInvalidStartTimeout    = "invalid run config start timeout"
	logMsgInvalidStopTimeout     = "invalid run config stop timeout"
	logMsgNilFactory             = "run factory is nil"
	logMsgBootstrapStartFailed   = "bootstrap start failed"
	logMsgBootstrapStartOK       = "bootstrap start completed"
	logMsgBootstrapStopFailed    = "bootstrap stop failed"
	logMsgBootstrapStopCompleted = "bootstrap stop completed"

	logFieldDurationMS   = "duration_ms"
	logFieldError        = "error"
	logFieldStartTimeout = "start_timeout"
	logFieldStopTimeout  = "stop_timeout"
	logFieldComponent    = "component"
	logFieldPhase        = "phase"

	logValueComponentBootstrap = "cmd_api_bootstrap"
	logValuePhaseConfig        = "config"
	logValuePhaseStart         = "start"
	logValuePhaseStop          = "stop"
)
