// Package token is a JWT-based implementation of output.Provider.
package token

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
)

const ExpTimeToken = 24 * time.Hour

// Provider implements output.Provider using HMAC-SHA256.
type Provider struct {
	secretKey string
}

// New builds a JWT token with the given secret key.
func New(secretKey string) *Provider {
	return &Provider{secretKey: secretKey}
}

// Generate creates a signed JWT with userID and expiration.
func (p *Provider) Generate(_ context.Context, userID uint64) (string, error) {
	claims := jwt.MapClaims{
		claimskeys.UserID: userID,
		claimskeys.Exp:    time.Now().Add(ExpTimeToken).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(p.secretKey))
}

// Verify checks signature/exp and returns the claims map on success.
func (p *Provider) Verify(_ context.Context, token string) (map[string]any, error) {
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
