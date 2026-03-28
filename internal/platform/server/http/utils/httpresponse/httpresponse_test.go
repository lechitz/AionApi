package httpresponse_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httperrors "github.com/lechitz/aion-api/internal/platform/server/http/errors"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/tracingkeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	sdkexport "go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type testTracer struct {
	tp       *tracesdk.TracerProvider
	exporter *sdkexport.InMemoryExporter
}

func newTestTracer(t *testing.T) *testTracer {
	t.Helper()
	exporter := sdkexport.NewInMemoryExporter()
	sp := tracesdk.NewSimpleSpanProcessor(exporter)
	tp := tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(sp))
	return &testTracer{tp: tp, exporter: exporter}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		payload         any
		headers         map[string]string
		wantContentType bool
		wantBody        bool
	}{
		{
			name:            "no_content",
			status:          http.StatusNoContent,
			payload:         map[string]string{"ok": "true"},
			wantContentType: false,
			wantBody:        false,
		},
		{
			name:            "ok_with_headers",
			status:          http.StatusCreated,
			payload:         map[string]string{"ok": "true"},
			headers:         map[string]string{"X-Test": "value"},
			wantContentType: true,
			wantBody:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			httpresponse.WriteJSON(rec, tt.status, tt.payload, tt.headers)

			assert.Equal(t, tt.status, rec.Code)
			if tt.headers != nil {
				assert.Equal(t, tt.headers["X-Test"], rec.Header().Get("X-Test"))
			}
			if tt.wantContentType {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			} else {
				assert.Empty(t, rec.Header().Get("Content-Type"))
			}

			if tt.wantBody {
				var got map[string]string
				require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
				assert.Equal(t, tt.payload, got)
			} else {
				assert.Empty(t, rec.Body.String())
			}
		})
	}
}

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	result := map[string]string{"id": "1"}
	httpresponse.WriteSuccess(rec, http.StatusCreated, result, "created")

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body httpresponse.ResponseBody
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, http.StatusCreated, body.Code)
	assert.Equal(t, "created", body.Message)
	assert.False(t, body.Date.IsZero())

	gotResult, ok := body.Result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "1", gotResult["id"])
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		message    string
		wantStatus int
	}{
		{
			name:       "bad_request_from_parse_user_id",
			err:        sharederrors.ErrParseUserID,
			message:    "invalid user id",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not_found_from_http_error",
			err:        httperrors.ErrResourceNotFound,
			message:    "missing",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			httpresponse.WriteError(rec, tt.err, tt.message, nil)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var body httpresponse.ResponseBody
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
			assert.Equal(t, tt.wantStatus, body.Code)
			assert.Equal(t, tt.message, body.Error)
			assert.Empty(t, body.Details)
			assert.False(t, body.Date.IsZero())
		})
	}
}

func TestWriteAuthAndDecodeError(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(http.ResponseWriter, error)
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name: "decode_error",
			fn: func(w http.ResponseWriter, err error) {
				httpresponse.WriteDecodeError(w, err, nil)
			},
			err:        sharederrors.NewValidationError("field", "bad"),
			wantStatus: http.StatusBadRequest,
			wantMsg:    "Invalid request body",
		},
		{
			name: "auth_error",
			fn: func(w http.ResponseWriter, err error) {
				httpresponse.WriteAuthError(w, err, nil)
			},
			err:        sharederrors.NewAuthenticationError("token"),
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tt.fn(rec, tt.err)

			assert.Equal(t, tt.wantStatus, rec.Code)
			var body httpresponse.ResponseBody
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
			assert.Equal(t, tt.wantStatus, body.Code)
			assert.Equal(t, tt.wantMsg, body.Error)
			assert.Empty(t, body.Details)
		})
	}
}

func TestWriteNoContent(t *testing.T) {
	rec := httptest.NewRecorder()
	httpresponse.WriteNoContent(rec, map[string]string{"X-Test": "value"})

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "value", rec.Header().Get("X-Test"))
	assert.Empty(t, rec.Body.String())
}

func TestSpanErrorResponses(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(ctx context.Context, w http.ResponseWriter, span trace.Span, err error)
		err        error
		message    string
		wantStatus int
	}{
		{
			name: "auth_error_span",
			fn: func(ctx context.Context, w http.ResponseWriter, span trace.Span, err error) {
				httpresponse.WriteAuthErrorSpan(ctx, w, span, err, nil)
			},
			err:        sharederrors.NewAuthenticationError("token"),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "decode_error_span",
			fn: func(ctx context.Context, w http.ResponseWriter, span trace.Span, err error) {
				httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, nil)
			},
			err:        sharederrors.NewValidationError("field", "bad"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation_error_span",
			fn: func(ctx context.Context, w http.ResponseWriter, span trace.Span, err error) {
				httpresponse.WriteValidationErrorSpan(ctx, w, span, err, nil)
			},
			err:        sharederrors.NewValidationError("field", "bad"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "domain_error_span",
			fn: func(ctx context.Context, w http.ResponseWriter, span trace.Span, err error) {
				httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "not found", nil)
			},
			err:        httperrors.ErrResourceNotFound,
			message:    "not found",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tTracer := newTestTracer(t)
			defer func() {
				_ = tTracer.tp.Shutdown(t.Context())
			}()

			tracer := tTracer.tp.Tracer("test")
			ctx, span := tracer.Start(t.Context(), "span")
			rec := httptest.NewRecorder()

			tt.fn(ctx, rec, span, tt.err)
			span.End()

			spans := tTracer.exporter.GetSpans()
			require.Len(t, spans, 1)
			spanStub := spans[0]
			assert.Equal(t, codes.Error, spanStub.Status.Code)

			status, ok := findIntAttribute(spanStub.Attributes, tracingkeys.HTTPStatusCodeKey)
			require.True(t, ok)
			assert.Equal(t, tt.wantStatus, status)

			assert.Equal(t, tt.wantStatus, rec.Code)
			var body httpresponse.ResponseBody
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
			assert.Equal(t, tt.wantStatus, body.Code)
			if tt.message != "" {
				assert.Equal(t, tt.message, body.Error)
			}
		})
	}
}

func findIntAttribute(attrs []attribute.KeyValue, key string) (int, bool) {
	for _, attr := range attrs {
		if string(attr.Key) != key {
			continue
		}
		if attr.Value.Type() != attribute.INT64 {
			return 0, false
		}
		return int(attr.Value.AsInt64()), true
	}
	return 0, false
}

func TestWriteValidationErrorSpan_UsesErrorMessage(t *testing.T) {
	tTracer := newTestTracer(t)
	defer func() {
		_ = tTracer.tp.Shutdown(t.Context())
	}()

	tracer := tTracer.tp.Tracer("test")
	ctx, span := tracer.Start(t.Context(), "span")
	rec := httptest.NewRecorder()

	err := sharederrors.NewValidationError("field", "required")
	httpresponse.WriteValidationErrorSpan(ctx, rec, span, err, nil)
	span.End()

	var body httpresponse.ResponseBody
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, err.Error(), body.Error)
}

func TestWriteError_WithCustomHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	err := errors.New("boom")
	httpresponse.WriteError(rec, err, "failed", nil, map[string]string{"X-Trace": "123"})

	assert.Equal(t, "123", rec.Header().Get("X-Trace"))
}
