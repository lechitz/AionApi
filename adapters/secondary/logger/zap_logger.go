// Package logger provides a wrapper around zap.SugaredLogger for structured logging.
package logger

import "go.uber.org/zap"

// ZapLoggerAdapter is a wrapper around zap.SugaredLogger to provide formatted logging methods.// ZapLoggerAdapter is a wrapper around zap.SugaredLogger for structured and formatted logging.
type ZapLoggerAdapter struct {
	sugar *zap.SugaredLogger
}

// NewZapLoggerAdapter creates a new instance of ZapLoggerAdapter wrapping a zap.SugaredLogger for structured logging.
func NewZapLoggerAdapter(sugar *zap.SugaredLogger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{sugar: sugar}
}

// Infof logs a formatted informational message using a format string and optional arguments.
func (l *ZapLoggerAdapter) Infof(format string, args ...any) {
	l.sugar.Infof(format, args...)
}

// Errorf logs an error message with a formatted string and additional// arguments Error.f logs an error message using a formatted string and optional arguments.
func (l *ZapLoggerAdapter) Errorf(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}

// Warnf logs a warning message using a formatted string and optional arguments.
func (l *ZapLoggerAdapter) Warnf(format string, args ...any) {
	l.sugar.Warnf(format, args...)
}

// Debugf logs a debug-level message using a formatted string and optional arguments.
func (l *ZapLoggerAdapter) Debugf(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}

// Infow logs an informational message and structured context key-value pairs. Parameters: msg is the message, keysAndValues are context data.
func (l *ZapLoggerAdapter) Infow(msg string, keysAndValues ...any) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Errorw logs an error message with a given string and additional key-value pairs for structured context.
func (l *ZapLoggerAdapter) Errorw(msg string, keysAndValues ...any) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Warnw logs a warning message with structured context using key-value pairs. Parameters: msg is the log message, keysAndValues are optional context data.
func (l *ZapLoggerAdapter) Warnw(msg string, keysAndValues ...any) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Debugw logs a debug-level message with structured context using key-value pairs. Parameters: msg is the log message, keysAndValues are context data.
func (l *ZapLoggerAdapter) Debugw(msg string, keysAndValues ...any) {
	l.sugar.Debugw(msg, keysAndValues...)
}
