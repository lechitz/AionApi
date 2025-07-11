// Package contextlogger provides a zap-based contextlogger implementation.
package contextlogger

import (
	"context"
	"log"
	"os"

	"github.com/lechitz/AionApi/internal/shared/ctxkeys"

	"github.com/lechitz/AionApi/internal/shared/contextutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	failedToFlushLogger = "Failed to flush context logger: %v"
)

// ZapLoggerContextual implements the ContextLogger interface (context-aware logging).
type ZapLoggerContextual struct {
	base *zap.SugaredLogger
}

// New initializes a zap.SugaredLogger and returns a ContextLogger and a cleanup function.
func New() (*ZapLoggerContextual, func()) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	infoWriter := zapcore.Lock(os.Stdout)
	errorWriter := zapcore.Lock(os.Stderr)

	infoCore := zapcore.NewCore(encoder, infoWriter, infoLevel)
	errorCore := zapcore.NewCore(encoder, errorWriter, errorLevel)

	tee := zapcore.NewTee(infoCore, errorCore)

	logger := zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar := logger.Sugar()

	cleanup := func() {
		if err := sugar.Sync(); err != nil {
			log.Printf(failedToFlushLogger, err)
		}
	}

	return &ZapLoggerContextual{base: sugar}, cleanup
}

// Infof logs a formatted info-level message.
func (l *ZapLoggerContextual) Infof(format string, args ...any) {
	l.base.Infof(format, args...)
}

// Errorf logs a formatted error-level message.
func (l *ZapLoggerContextual) Errorf(format string, args ...any) {
	l.base.Errorf(format, args...)
}

// Debugf logs a formatted debug-level message.
func (l *ZapLoggerContextual) Debugf(format string, args ...any) {
	l.base.Debugf(format, args...)
}

// Warnf logs a formatted warn-level message.
func (l *ZapLoggerContextual) Warnf(format string, args ...any) {
	l.base.Warnf(format, args...)
}

// Infow logs a structured info-level message.
func (l *ZapLoggerContextual) Infow(msg string, keysAndValues ...any) {
	l.base.Infow(msg, keysAndValues...)
}

// Errorw logs a structured error-level message.
func (l *ZapLoggerContextual) Errorw(msg string, keysAndValues ...any) {
	l.base.Errorw(msg, keysAndValues...)
}

// Debugw logs a structured debug-level message.
func (l *ZapLoggerContextual) Debugw(msg string, keysAndValues ...any) {
	l.base.Debugw(msg, keysAndValues...)
}

// Warnw logs a structured warn-level message.
func (l *ZapLoggerContextual) Warnw(msg string, keysAndValues ...any) {
	l.base.Warnw(msg, keysAndValues...)
}

// InfowCtx adds contextual data from context and logs a structured info-level message.
func (l *ZapLoggerContextual) InfowCtx(ctx context.Context, msg string, keysAndValues ...any) {
	fields := enrichFieldsFromContext(ctx)
	l.base.Infow(msg, append(fields, keysAndValues...)...)
}

// ErrorwCtx adds contextual data from context and logs a structured error-level message.
func (l *ZapLoggerContextual) ErrorwCtx(ctx context.Context, msg string, keysAndValues ...any) {
	fields := enrichFieldsFromContext(ctx)
	l.base.Errorw(msg, append(fields, keysAndValues...)...)
}

// DebugwCtx adds contextual data from context and logs a structured debug-level message.
func (l *ZapLoggerContextual) DebugwCtx(ctx context.Context, msg string, keysAndValues ...any) {
	fields := enrichFieldsFromContext(ctx)
	l.base.Debugw(msg, append(fields, keysAndValues...)...)
}

// WarnwCtx adds contextual data from context and logs a structured warn-level message.
func (l *ZapLoggerContextual) WarnwCtx(ctx context.Context, msg string, keysAndValues ...any) {
	fields := enrichFieldsFromContext(ctx)
	l.base.Warnw(msg, append(fields, keysAndValues...)...)
}

// enrichFieldsFromContext extracts relevant request-scoped fields (e.g., request_id, trace_id, user_id) from context for structured logging.
func enrichFieldsFromContext(ctx context.Context) []any {
	var fields []any
	if reqID := contextutils.GetRequestID(ctx); reqID != "" {
		fields = append(fields, ctxkeys.CtxKeyRequestID, reqID)
	}
	if traceID := contextutils.GetTraceID(ctx); traceID != "" {
		fields = append(fields, ctxkeys.CtxKeyTraceID, traceID)
	}
	if userID := contextutils.GetUserID(ctx); userID != "" {
		fields = append(fields, ctxkeys.UserID, userID)
	}
	return fields
}
