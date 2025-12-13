package httpclient_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/platform/httpclient"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestNewInstrumentedClient_InjectsTraceparent(t *testing.T) {
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(t.Context())
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Traceparent") == "" && r.Header.Get("traceparent") == "" {
			t.Errorf("expected traceparent header to be present")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := httpclient.NewInstrumentedClient(httpclient.Options{Timeout: 5 * time.Second})
	ctx := t.Context()

	tr := otel.Tracer("test")
	ctx, span := tr.Start(ctx, "unit-test")
	defer span.End()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}
