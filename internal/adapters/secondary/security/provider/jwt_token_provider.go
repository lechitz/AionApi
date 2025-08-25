// Package provider is a JWT-based implementation of output.TokenProvider.
package provider

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
)

const ExpTimeToken = 24 * time.Hour

// TokenProvider implements output.TokenProvider using HMAC-SHA256.
type TokenProvider struct {
	secretKey string
}

// New builds a JWT token provider with the given secret key.
func New(secretKey string) *TokenProvider {
	return &TokenProvider{secretKey: secretKey}
}

// Generate creates a signed JWT with userID and expiration.
func (p *TokenProvider) Generate(_ context.Context, userID uint64) (string, error) {
	claims := jwt.MapClaims{
		claimskeys.UserID: userID,
		claimskeys.Exp:    time.Now().Add(ExpTimeToken).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(p.secretKey))
}

// Verify checks signature/exp and returns the claims map on success.
func (p *TokenProvider) Verify(_ context.Context, token string) (map[string]any, error) {
	parsed, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return nil, err
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
