package contextlogger_test

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/adapter/secondary/contextlogger"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
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
