package httpclient_test

import (
	"context"
	"net"
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

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot start test listener: %v", err)
	}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Traceparent") == "" && r.Header.Get("traceparent") == "" {
			t.Errorf("expected traceparent header to be present")
		}
		w.WriteHeader(http.StatusOK)
	}))
	server.Listener = listener
	server.Start()
	defer server.Close()

	client := httpclient.NewInstrumentedClient(httpclient.Options{Timeout: 5 * time.Second})
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
