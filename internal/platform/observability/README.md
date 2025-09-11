# Observability (Platform · OpenTelemetry)

**Folder:** `internal/platform/observability`
**Subpackages:**

* `tracer/` — tracing pipeline (OTLP over HTTP)
* `metric/` — metrics pipeline (OTLP over HTTP)
* `helpers.go` — small utils (e.g., OTEL header parsing)

## Responsibility

* Provide a **single place** to bootstrap OpenTelemetry **traces** and **metrics** for the whole platform.
* Expose helpers to **configure exporters** (endpoint, headers, security) from `config.ObservabilityConfig`.
* Set **resource attributes** (service name, version, environment) and standardize logging of startup/shutdown.

---

## What it sets up

* **TracerProvider** (OTLP/HTTP) — spans for HTTP handlers, use cases, repositories, etc.
* **MeterProvider** (OTLP/HTTP) — counters, histograms, and gauges.
* **Common Resource** — `service.name`, `service.version`, `deployment.environment`, and other attributes used across both pipelines.
* **Headers parsing** — `ParseHeaders("k=v,k2=v2") → map[string]string` for custom OTLP headers.

---

## How it works

* Both `tracer/` and `metric/` read from `config.ObservabilityConfig` and a `logger.ContextLogger`.
* They build an OTLP **HTTP** exporter, attach the shared **Resource**, set the global provider via `otel.SetTracerProvider`/`otel.SetMeterProvider`, and return a **shutdown function** to flush on exit.
* Exporters support:

    * `OTEL_EXPORTER_OTLP_ENDPOINT` (HTTP)
    * Optional headers (CSV → map via `ParseHeaders`)
    * Insecure toggle (dev) and simple sampler knobs

> See `tracer/tracing.go` and `metric/metrics.go` for the concrete initialization functions.

---

## Configuration (from `config.ObservabilityConfig`)

Typical fields (names may vary; see `internal/platform/config`):

* **Endpoint:** `OTEL_EXPORTER_OTLP_ENDPOINT` (e.g., `http://localhost:4318`)
* **Headers (CSV):** `OTEL_EXPORTER_OTLP_HEADERS="x-api-key=abc,another=val"`
* **Insecure:** allow plain HTTP for local/dev
* **Service metadata:** `OTEL_SERVICE_NAME`, `OTEL_SERVICE_VERSION`
* **Sampler:** parent/always/ratio (e.g., `0.1` in dev)

---

## Usage

```go
// Pseudocode – check actual function names in tracer/ and metric/
shutdowns := make([]func(context.Context) error, 0, 2)

// Tracing
if sd, err := tracer.Setup(ctx, cfg.Observability, log); err != nil {
  log.Errorw("otel.tracer.setup.failed", "error", err)
} else {
  shutdowns = append(shutdowns, sd)
}

// Metrics
if sd, err := metric.Setup(ctx, cfg.Observability, log); err != nil {
  log.Errorw("otel.metric.setup.failed", "error", err)
} else {
  shutdowns = append(shutdowns, sd)
}

// On graceful shutdown:
for i := len(shutdowns) - 1; i >= 0; i-- {
  _ = shutdowns[i](context.Background())
}
```

Then, anywhere in the code:

```go
tr := otel.Tracer("aionapi.category.repo")
ctx, span := tr.Start(ctx, "repo.list_all")
defer span.End()
```

---

## Design notes

* **HTTP exporters only** (collector-friendly, firewall-friendly).
* **No business logic** here; keep this package pure platform wiring.
* **Global providers** are set once at bootstrap; downstream code just calls `otel.Tracer(...)` / `otel.Meter(...)`.
* Reuse **common attributes/keys** from `internal/shared/constants` to keep telemetry consistent.

---

## Testing hints

* Local dev: run an OTEL Collector with HTTP receivers (`:4318`).
* To avoid noisy tests, **skip initialization** (e.g., empty endpoint) or point to a test collector.
* For header parsing unit tests, use `ParseHeaders("a=b,c=d")`.

---

## Troubleshooting

* **No data arriving:** check endpoint (`http://host:4318`), headers, and network access.
* **Empty service labels:** ensure `OTEL_SERVICE_NAME` and `OTEL_SERVICE_VERSION` are set in env/config.
* **TLS errors:** use the insecure flag for local dev or configure proper HTTPS on the collector.
