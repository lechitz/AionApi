package requestid_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lechitz/aion-api/internal/platform/server/http/middleware/requestid"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
)

func TestRequestIDMiddleware_GeneratesAndInjectsWhenMissing(t *testing.T) {
	mw := requestid.New()

	var ctxReqID string
	h := mw(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value(ctxkeys.RequestID).(string)
		ctxReqID = val
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	headerReqID := w.Result().Header.Get(commonkeys.XRequestID)
	if headerReqID == "" {
		t.Fatalf("expected response request id header")
	}
	if _, err := uuid.Parse(headerReqID); err != nil {
		t.Fatalf("expected UUID request id, got %q", headerReqID)
	}
	if ctxReqID != headerReqID {
		t.Fatalf("expected context request id to match header: %q != %q", ctxReqID, headerReqID)
	}
}

func TestRequestIDMiddleware_PreservesValidUUIDHeader(t *testing.T) {
	mw := requestid.New()
	valid := uuid.NewString()

	var ctxReqID string
	h := mw(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value(ctxkeys.RequestID).(string)
		ctxReqID = val
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set(commonkeys.XRequestID, valid)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	got := w.Result().Header.Get(commonkeys.XRequestID)
	if got != valid || ctxReqID != valid {
		t.Fatalf("expected valid request id to be preserved")
	}
}

func TestRequestIDMiddleware_ReplacesInvalidOrTooLongHeader(t *testing.T) {
	mw := requestid.New()
	tooLongInvalid := strings.Repeat("x", 129)

	h := mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set(commonkeys.XRequestID, tooLongInvalid)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	got := w.Result().Header.Get(commonkeys.XRequestID)
	if got == "" || got == tooLongInvalid {
		t.Fatalf("expected invalid request id to be replaced")
	}
	if _, err := uuid.Parse(got); err != nil {
		t.Fatalf("expected generated UUID after replacement, got %q", got)
	}
}
