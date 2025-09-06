// Package logger output contextlogger provides an interface for structured and leveled logging within applications.
package logger

import "context"

// ContextLogger is an interface for structured and leveled logging within applications.
type ContextLogger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
	Warnf(format string, args ...any)

	Infow(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)

	InfowCtx(ctx context.Context, msg string, keysAndValues ...any)
	ErrorwCtx(ctx context.Context, msg string, keysAndValues ...any)
	WarnwCtx(ctx context.Context, msg string, keysAndValues ...any)
	DebugwCtx(ctx context.Context, msg string, keysAndValues ...any)
}
