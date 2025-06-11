// Package logger provides an interface for structured and leveled logging within applications.
package logger

// Logger is an interface for structured and leveled logging within applications.
// Infof logs informational messages with a formatted string.
// Errorf logs error messages with a formatted string.
// Debugf logs debug-level messages with a formatted string.
// Warnf logs warning messages with a formatted string.
// Infow logs informational messages with structured key-value pairs.
// Errorw logs error messages with structured key-value pairs.
// Debugw logs debug-level messages with structured key-value pairs.
// Warnw logs warning messages with structured key-value pairs.
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
	Warnf(format string, args ...any)

	Infow(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
}
