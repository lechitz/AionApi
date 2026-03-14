package recovery_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	ghandler "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
	"github.com/lechitz/AionApi/internal/platform/server/http/middleware/recovery"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/lechitz/AionApi/tests/setup"
	"go.uber.org/mock/gomock"
)

func TestRecoveryMiddleware_RecoversPanicAndReturns500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockContextLogger(ctrl)
	setup.ExpectLoggerDefaultBehavior(log)
	log.EXPECT().
		ErrorwCtx(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
		).
		AnyTimes()
	log.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Warnw(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Debugw(gomock.Any(), gomock.Any()).AnyTimes()

	recoveryHandler := ghandler.New(log, config.GeneralConfig{})
	mw := recovery.New(recoveryHandler)

	h := mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 after panic recovery, got %d", w.Code)
	}
}

func TestRecoveryMiddleware_PassesThroughWithoutPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockContextLogger(ctrl)
	setup.ExpectLoggerDefaultBehavior(log)
	log.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Errorw(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Warnw(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Debugw(gomock.Any(), gomock.Any()).AnyTimes()

	recoveryHandler := ghandler.New(log, config.GeneralConfig{})
	mw := recovery.New(recoveryHandler)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/ok", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected pass-through status, got %d", w.Code)
	}
}
