package contextlogger_test

import (
	"context"
	"testing"

	"github.com/lechitz/aion-api/internal/adapter/secondary/contextlogger"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestContextHelpers_Getters(t *testing.T) {
	ctx := t.Context()
	ctx = context.WithValue(ctx, ctxkeys.RequestID, []byte("req-1"))
	ctx = context.WithValue(ctx, ctxkeys.TraceID, int64(12345))
	ctx = context.WithValue(ctx, ctxkeys.UserID, uint64(99))

	require.Equal(t, "req-1", contextlogger.GetRequestID(ctx))
	require.Equal(t, "12345", contextlogger.GetTraceID(ctx))
	require.Equal(t, "99", contextlogger.GetUserID(ctx))
}

func TestEnrichFieldsFromContext_UsesSpanContextWhenValid(t *testing.T) {
	ctx := t.Context()
	ctx = context.WithValue(ctx, ctxkeys.RequestID, "req-2")
	ctx = context.WithValue(ctx, ctxkeys.TraceID, "fallback-trace")
	ctx = context.WithValue(ctx, ctxkeys.UserID, "user-2")

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:     trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})
	ctx = trace.ContextWithSpanContext(ctx, sc)

	fields := contextlogger.EnrichFieldsFromContext(ctx)
	m := fieldsToMap(fields)

	require.Equal(t, "req-2", m[string(ctxkeys.RequestID)])
	require.Equal(t, sc.TraceID().String(), m[string(ctxkeys.TraceID)])
	require.Equal(t, sc.SpanID().String(), m[string(ctxkeys.SpanID)])
	require.Equal(t, "user-2", m[string(ctxkeys.UserID)])
}

func TestEnrichFieldsFromContext_FallbackTraceID(t *testing.T) {
	ctx := t.Context()
	ctx = context.WithValue(ctx, ctxkeys.TraceID, "trace-fallback")

	fields := contextlogger.EnrichFieldsFromContext(ctx)
	m := fieldsToMap(fields)
	require.Equal(t, "trace-fallback", m[string(ctxkeys.TraceID)])
}

func TestNewLoggerAndCleanup(t *testing.T) {
	logger, cleanup := contextlogger.New()
	require.NotNil(t, logger)
	require.NotNil(t, cleanup)
	cleanup()
}

func TestLoggerMethods_DoNotPanic(t *testing.T) {
	logger, cleanup := contextlogger.New()
	t.Cleanup(cleanup)

	ctx := context.WithValue(t.Context(), ctxkeys.RequestID, 123)
	ctx = context.WithValue(ctx, ctxkeys.TraceID, []byte("trace-b"))
	ctx = context.WithValue(ctx, ctxkeys.UserID, int64(7))

	logger.Infof("hello %s", "world")
	logger.Errorf("err %d", 1)
	logger.Debugf("dbg %t", true)
	logger.Warnf("warn %s", "x")

	logger.Infow("msg", "k1", "v1")
	logger.Errorw("msg", "k1", "v1")
	logger.Debugw("msg", "k1", "v1")
	logger.Warnw("msg", "k1", "v1")

	logger.InfowCtx(ctx, "ctx-msg", "k2", "v2")
	logger.ErrorwCtx(ctx, "ctx-msg", "k2", "v2")
	logger.DebugwCtx(ctx, "ctx-msg", "k2", "v2")
	logger.WarnwCtx(ctx, "ctx-msg", "k2", "v2")

	require.Equal(t, "123", contextlogger.GetRequestID(ctx))
	require.Equal(t, "trace-b", contextlogger.GetTraceID(ctx))
	require.Equal(t, "7", contextlogger.GetUserID(ctx))
}

func fieldsToMap(fields []any) map[string]string {
	out := make(map[string]string)
	for i := 0; i+1 < len(fields); i += 2 {
		k, ok := fields[i].(string)
		if !ok {
			continue
		}
		v, ok := fields[i+1].(string)
		if !ok {
			continue
		}
		out[k] = v
	}
	return out
}
