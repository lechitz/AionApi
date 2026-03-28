package domain_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/auth/core/domain"
	"github.com/stretchr/testify/require"
)

func TestTokenConstructorsAndConversions(t *testing.T) {
	access := domain.NewAccessToken("access-token", 42)
	refresh := domain.NewRefreshToken("refresh-token", 99)

	require.Equal(t, domain.TokenTypeAccess, access.Type)
	require.Equal(t, domain.TokenTypeRefresh, refresh.Type)

	accessAuth := access.ToAuth()
	refreshAuth := refresh.ToAuth()
	require.Equal(t, domain.Auth{Token: "access-token", Key: 42, Type: domain.TokenTypeAccess}, accessAuth)
	require.Equal(t, domain.Auth{Token: "refresh-token", Key: 99, Type: domain.TokenTypeRefresh}, refreshAuth)

	convertedAccess := domain.AccessTokenFromAuth(domain.Auth{Token: "a", Key: 7, Type: domain.TokenTypeAccess})
	convertedRefresh := domain.RefreshTokenFromAuth(domain.Auth{Token: "r", Key: 8, Type: domain.TokenTypeRefresh})
	require.Equal(t, domain.AccessToken{Token: "a", Key: 7, Type: domain.TokenTypeAccess}, convertedAccess)
	require.Equal(t, domain.RefreshToken{Token: "r", Key: 8, Type: domain.TokenTypeRefresh}, convertedRefresh)
}
