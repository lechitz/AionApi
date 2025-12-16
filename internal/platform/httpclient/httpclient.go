package httpclient

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Options configures the instrumented HTTP client.
type Options struct {
	Timeout                time.Duration
	BaseTransport          http.RoundTripper
	DefaultHeaders         map[string]string
	DisableInstrumentation bool
	OtelOptions            []otelhttp.Option
}

// NewInstrumentedClient returns an *http.Client which uses an instrumented transport
// (otelhttp.Transport) when instrumentation is enabled. The returned client is safe
// for concurrent use.
func NewInstrumentedClient(opts Options) *http.Client {
	transport := opts.BaseTransport
	if transport == nil {
		transport = http.DefaultTransport
	}

	var rt http.RoundTripper
	if opts.DisableInstrumentation {
		rt = transport
	} else {
		// wrap transport with otelhttp to create spans and inject propagation headers
		rt = otelhttp.NewTransport(transport, opts.OtelOptions...)
	}

	c := &http.Client{
		Transport: rt,
	}

	if opts.Timeout > 0 {
		c.Timeout = opts.Timeout
	} else {
		c.Timeout = 15 * time.Second
	}

	return c
}
