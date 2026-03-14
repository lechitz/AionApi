package handler_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/lechitz/AionApi/internal/chat/adapter/primary/http/handler"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/stretchr/testify/require"
)

func newRestrictedChatServer(t *testing.T, fn http.HandlerFunc) *httptest.Server {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot start test listener: %v", err)
	}

	srv := httptest.NewUnstartedServer(fn)
	srv.Listener = listener
	srv.Start()
	t.Cleanup(srv.Close)
	return srv
}

func newVoiceRequest(t *testing.T, audio []byte, language string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if audio != nil {
		part, err := writer.CreateFormFile("audio", "sample.wav")
		require.NoError(t, err)
		_, err = part.Write(audio)
		require.NoError(t, err)
	}

	if language != "" {
		require.NoError(t, writer.WriteField("language", language))
	}

	require.NoError(t, writer.Close())

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chat/audio", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestChatVoice_Success(t *testing.T) {
	chatSrv := newRestrictedChatServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/process-audio" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		r.Body = http.MaxBytesReader(w, r.Body, handler.MaxAudioSize)
		if err := r.ParseMultipartForm(handler.MaxAudioSize); err != nil {
			t.Errorf("unexpected multipart error: %v", err)
		}
		if r.FormValue("user_id") != "7" {
			t.Errorf("unexpected user_id: %s", r.FormValue("user_id"))
		}
		if r.FormValue("language") != "pt" {
			t.Errorf("unexpected language: %s", r.FormValue("language"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"ok"}`))
	})

	h := handler.New(mockChatService{}, &config.Config{
		AionChat: config.AionChatConfig{BaseURL: chatSrv.URL},
	}, mockLogger{})

	req := newVoiceRequest(t, []byte("audio-bytes"), "pt")
	req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(7)))
	rec := httptest.NewRecorder()

	h.ChatVoice(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	require.Contains(t, rec.Body.String(), `"message":"ok"`)
}

func TestChatVoice_AuthErrors(t *testing.T) {
	h := handler.New(mockChatService{}, &config.Config{
		AionChat: config.AionChatConfig{BaseURL: "http://127.0.0.1:1"},
	}, mockLogger{})

	t.Run("missing user id", func(t *testing.T) {
		req := newVoiceRequest(t, []byte("audio"), "")
		rec := httptest.NewRecorder()
		h.ChatVoice(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid user id type", func(t *testing.T) {
		req := newVoiceRequest(t, []byte("audio"), "")
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, "7"))
		rec := httptest.NewRecorder()
		h.ChatVoice(rec, req)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestChatVoice_RequestValidation(t *testing.T) {
	h := handler.New(mockChatService{}, &config.Config{
		AionChat: config.AionChatConfig{BaseURL: "http://127.0.0.1:1"},
	}, mockLogger{})

	t.Run("missing audio file", func(t *testing.T) {
		req := newVoiceRequest(t, nil, "pt")
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(1)))
		rec := httptest.NewRecorder()
		h.ChatVoice(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("audio file too large", func(t *testing.T) {
		tooLargeAudio := bytes.Repeat([]byte("a"), handler.MaxAudioSize+1)
		req := newVoiceRequest(t, tooLargeAudio, "")
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(1)))
		rec := httptest.NewRecorder()
		h.ChatVoice(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "too large")
	})
}

func TestChatVoice_ForwardingErrorsAndServiceStatus(t *testing.T) {
	t.Run("invalid upstream URL", func(t *testing.T) {
		h := handler.New(mockChatService{}, &config.Config{
			AionChat: config.AionChatConfig{BaseURL: "://bad-url"},
		}, mockLogger{})

		req := newVoiceRequest(t, []byte("audio"), "")
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(1)))
		rec := httptest.NewRecorder()
		h.ChatVoice(rec, req)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("upstream returns non-200", func(t *testing.T) {
		chatSrv := newRestrictedChatServer(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"unavailable"}`))
		})
		h := handler.New(mockChatService{}, &config.Config{
			AionChat: config.AionChatConfig{BaseURL: chatSrv.URL},
		}, mockLogger{})

		req := newVoiceRequest(t, []byte("audio"), "")
		req = req.WithContext(context.WithValue(t.Context(), ctxkeys.UserID, uint64(1)))
		rec := httptest.NewRecorder()
		h.ChatVoice(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)
		require.Contains(t, rec.Body.String(), "unavailable")
	})
}
