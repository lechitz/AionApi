package token_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/adapter/secondary/token"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/stretchr/testify/require"
)

func TestProvider_GenerateAccessAndVerify(t *testing.T) {
	p := token.NewProvider("test-secret")

	signed, err := p.GenerateAccessToken(42, map[string]any{"role": "admin"})
	require.NoError(t, err)
	require.NotEmpty(t, signed)

	claims, err := p.Verify(signed)
	require.NoError(t, err)
	require.Equal(t, "42", claims[claimskeys.UserID])
	require.Equal(t, "admin", claims["role"])
}

func TestProvider_GenerateRefreshAndVerify(t *testing.T) {
	p := token.NewProvider("test-secret")

	signed, err := p.GenerateRefreshToken(77)
	require.NoError(t, err)

	claims, err := p.Verify(signed)
	require.NoError(t, err)
	require.Equal(t, "77", claims[claimskeys.UserID])
}

func TestProvider_VerifyErrors(t *testing.T) {
	p := token.NewProvider("test-secret")
	other := token.NewProvider("other-secret")

	signed, err := other.GenerateAccessToken(1, nil)
	require.NoError(t, err)

	_, err = p.Verify("not-a-token")
	require.Error(t, err)

	_, err = p.Verify(signed)
	require.Error(t, err)
}

func TestProvider_VerifyExpiredToken(t *testing.T) {
	p := token.NewProvider("test-secret")

	claims := jwt.MapClaims{
		claimskeys.UserID: "9",
		claimskeys.Exp:    time.Now().Add(-time.Minute).Unix(),
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tk.SignedString([]byte("test-secret"))
	require.NoError(t, err)

	_, err = p.Verify(signed)
	require.Error(t, err)
}
