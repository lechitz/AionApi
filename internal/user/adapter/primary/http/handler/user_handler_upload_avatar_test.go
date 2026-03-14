package handler_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	handler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
	"github.com/stretchr/testify/require"
)

func buildAvatarMultipart(t *testing.T, fieldName, filename string, payload []byte) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)
	_, err = part.Write(payload)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return body, writer.FormDataContentType()
}

func TestUploadAvatar_Success(t *testing.T) {
	svc := &mockUserService{
		uploadAvatarFn: func(_ context.Context, _ input.UploadAvatarCommand) (string, string, int64, error) {
			return "data:image/png;base64,AA==", "image/png", 2, nil
		},
	}
	h := handler.New(svc, &config.Config{}, mockLogger{})

	body, contentType := buildAvatarMultipart(t, "avatar", "avatar.png", []byte{0x89, 0x50})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/user/avatar/upload", body)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	h.UploadAvatar(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUploadAvatar_MissingFile(t *testing.T) {
	svc := &mockUserService{}
	h := handler.New(svc, &config.Config{}, mockLogger{})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.Close())

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/user/avatar/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	h.UploadAvatar(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadAvatar_UsecaseError(t *testing.T) {
	svc := &mockUserService{
		uploadAvatarFn: func(_ context.Context, _ input.UploadAvatarCommand) (string, string, int64, error) {
			return "", "", 0, context.Canceled
		},
	}
	h := handler.New(svc, &config.Config{}, mockLogger{})

	body, contentType := buildAvatarMultipart(t, "avatar", "avatar.png", []byte{0x89, 0x50})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/user/avatar/upload", body)
	req.Header.Set("Content-Type", contentType)
	rec := httptest.NewRecorder()

	h.UploadAvatar(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
