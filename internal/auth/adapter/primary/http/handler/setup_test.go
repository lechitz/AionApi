package handler_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	sdkexport "go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// TestTracer holds a tracer provider and in-memory exporter useful for assertions in tests.
type TestTracer struct {
	TP       *tracesdk.TracerProvider
	Exporter *sdkexport.InMemoryExporter
}

// NewTestTracer creates a new in-memory tracer provider and sets it as the global provider.
func NewTestTracer(t *testing.T) *TestTracer {
	t.Helper()
	exp := sdkexport.NewInMemoryExporter()
	sp := tracesdk.NewSimpleSpanProcessor(exp)
	tp := tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(sp))
	otel.SetTracerProvider(tp)
	return &TestTracer{TP: tp, Exporter: exp}
}

// Shutdown shuts down the tracer provider.
func (tt *TestTracer) Shutdown(ctx context.Context) error {
	return tt.TP.Shutdown(ctx)
}

// Spans returns recorded spans since exporter was created.
func (tt *TestTracer) Spans() sdkexport.SpanStubs {
	return tt.Exporter.GetSpans()
}

// Clear exported spans (useful between assertions).
func (tt *TestTracer) Reset() {
	tt.Exporter.Reset()
}
