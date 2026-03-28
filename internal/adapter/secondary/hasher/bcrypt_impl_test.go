package hasher_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/adapter/secondary/hasher"
	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_HashAndCompare(t *testing.T) {
	h := hasher.New()

	hashed, err := h.Hash("my-secret")
	require.NoError(t, err)
	require.NotEmpty(t, hashed)
	require.NotEqual(t, "my-secret", hashed)

	require.NoError(t, h.Compare(hashed, "my-secret"))
	require.Error(t, h.Compare(hashed, "wrong-secret"))
}
