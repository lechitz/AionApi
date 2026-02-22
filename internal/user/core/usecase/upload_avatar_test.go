package usecase_test

import (
	"bytes"
	"context"
	"testing"

	authOutput "github.com/lechitz/AionApi/internal/auth/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/hasher"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
	userOutput "github.com/lechitz/AionApi/internal/user/core/ports/output"
	"github.com/lechitz/AionApi/internal/user/core/usecase"
	"github.com/stretchr/testify/require"
)

type (
	noopUserRepo     struct{ userOutput.UserRepository }
	noopUserCache    struct{ userOutput.UserCache }
	noopAuthStore    struct{ authOutput.AuthStore }
	noopAuthProvider struct{ authOutput.AuthProvider }
	noopHasher       struct{ hasher.Hasher }
	noopLogger       struct{ logger.ContextLogger }
)

type mockAvatarStorage struct {
	url string
	err error
}

type fakeMultipartFile struct {
	*bytes.Reader
}

func (f *fakeMultipartFile) Close() error { return nil }

func (m *mockAvatarStorage) UploadAvatar(_ context.Context, _ userOutput.AvatarUploadInput) (string, error) {
	return m.url, m.err
}

func TestUploadAvatar_Success(t *testing.T) {
	svc := usecase.NewService(
		noopUserRepo{},
		nil,
		noopUserCache{},
		&mockAvatarStorage{url: "http://localhost:4566/aion-assets/avatars/a.png"},
		noopAuthStore{},
		noopAuthProvider{},
		noopHasher{},
		noopLogger{},
	)

	content := append([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, bytes.Repeat([]byte{0}, 200)...)
	url, contentType, size, err := svc.UploadAvatar(t.Context(), input.UploadAvatarCommand{
		File:        &fakeMultipartFile{Reader: bytes.NewReader(content)},
		Filename:    "avatar.png",
		SizeBytes:   int64(len(content)),
		ContentType: "image/png",
	})
	require.NoError(t, err)
	require.Equal(t, "http://localhost:4566/aion-assets/avatars/a.png", url)
	require.Equal(t, "image/png", contentType)
	require.Equal(t, int64(len(content)), size)
}

func TestUploadAvatar_TooLarge(t *testing.T) {
	svc := usecase.NewService(
		noopUserRepo{},
		nil,
		noopUserCache{},
		&mockAvatarStorage{url: "unused"},
		noopAuthStore{},
		noopAuthProvider{},
		noopHasher{},
		noopLogger{},
	)

	tooLarge := bytes.Repeat([]byte{1}, (20<<20)+1)
	_, _, _, err := svc.UploadAvatar(t.Context(), input.UploadAvatarCommand{
		File:        &fakeMultipartFile{Reader: bytes.NewReader(tooLarge)},
		Filename:    "large.png",
		SizeBytes:   int64(len(tooLarge)),
		ContentType: "image/png",
	})
	require.Error(t, err)
}
