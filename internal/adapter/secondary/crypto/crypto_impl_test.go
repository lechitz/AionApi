package crypto_test

import (
	"encoding/base64"
	"testing"

	"github.com/lechitz/aion-api/internal/adapter/secondary/crypto"
	"github.com/stretchr/testify/require"
)

func TestKeyGenerator_Generate(t *testing.T) {
	g := crypto.New()

	key, err := g.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, key)

	raw, err := base64.StdEncoding.DecodeString(key)
	require.NoError(t, err)
	require.Len(t, raw, crypto.KeyLength)
}
