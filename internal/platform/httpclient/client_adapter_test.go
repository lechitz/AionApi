package httpclient_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/httpclient"
)

func newRestrictedSafeServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot start test listener: %v", err)
	}

	server := httptest.NewUnstartedServer(handler)
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)
	return server
}

func TestNewClient_NilUsesDefaultClient(t *testing.T) {
	c := httpclient.NewClient(nil)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestClientAdapter_Get_Success(t *testing.T) {
	srv := newRestrictedSafeServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	c := httpclient.NewClient(&http.Client{})
	resp, err := c.Get(t.Context(), srv.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestClientAdapter_Get_InvalidURL(t *testing.T) {
	c := httpclient.NewClient(&http.Client{})
	resp, err := c.Get(t.Context(), "://bad-url")
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestClientAdapter_Post_Success(t *testing.T) {
	srv := newRestrictedSafeServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content-type: %s", got)
		}
		w.WriteHeader(http.StatusCreated)
	})

	c := httpclient.NewClient(&http.Client{})
	resp, err := c.Post(t.Context(), srv.URL, "application/json", map[string]any{"ok": true})
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestClientAdapter_Post_NilBody_Success(t *testing.T) {
	srv := newRestrictedSafeServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	c := httpclient.NewClient(&http.Client{})
	resp, err := c.Post(t.Context(), srv.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("Post with nil body failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestClientAdapter_Post_MarshalError(t *testing.T) {
	c := httpclient.NewClient(&http.Client{})
	resp, err := c.Post(t.Context(), "http://example.com", "application/json", func() {})
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestClientAdapter_Post_InvalidURL(t *testing.T) {
	c := httpclient.NewClient(&http.Client{})
	resp, err := c.Post(t.Context(), "://bad-url", "application/json", map[string]any{"ok": true})
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err == nil {
		t.Fatal("expected request build error")
	}
}

func TestNewInstrumentedClient_DisableInstrumentation_NoTraceparentInjection(t *testing.T) {
	var seenTraceHeader string
	srv := newRestrictedSafeServer(t, func(w http.ResponseWriter, r *http.Request) {
		seenTraceHeader = r.Header.Get("Traceparent")
		if seenTraceHeader == "" {
			seenTraceHeader = r.Header.Get("traceparent")
		}
		w.WriteHeader(http.StatusOK)
	})

	client := httpclient.NewInstrumentedClient(httpclient.Options{
		DisableInstrumentation: true,
	})

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if strings.TrimSpace(seenTraceHeader) != "" {
		t.Fatalf("expected no traceparent with instrumentation disabled, got %q", seenTraceHeader)
	}
}
