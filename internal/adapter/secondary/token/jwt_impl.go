// Package token is a JWT-based implementation of output.Provider.
package token

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
)

// ExpTimeToken is the default expiration time for JWT tokens.
const ExpTimeToken = 24 * time.Hour

// Provider implements output.Provider using HMAC-SHA256.
type Provider struct {
	secretKey string
}

// NewProvider builds a JWT token with the given secret key.
func NewProvider(secretKey string) *Provider {
	return &Provider{secretKey: secretKey}
}

// Generate creates a signed JWT with userID and expiration.
func (p *Provider) Generate(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		claimskeys.UserID: strconv.FormatUint(userID, 10),
		claimskeys.Exp:    time.Now().Add(ExpTimeToken).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(p.secretKey))
}

// Verify checks signature and expiration, returning the claims map on success.
func (p *Provider) Verify(token string) (map[string]any, error) {
	claims := jwt.MapClaims{}

	parsedToken, err := jwt.ParseWithClaims(
		token,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(p.secretKey), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if !expOK(claims[claimskeys.Exp]) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// GenerateWithClaims creates a signed JWT with userID, expiration and extra claims (e.g., roles).
func (p *Provider) GenerateWithClaims(userID uint64, extra map[string]any) (string, error) {
	claims := jwt.MapClaims{
		claimskeys.UserID: strconv.FormatUint(userID, 10),
		claimskeys.Exp:    time.Now().Add(ExpTimeToken).Unix(),
	}
	for k, v := range extra {
		claims[k] = v
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(p.secretKey))
}

// expOK checks if an "exp" value (various JSON-decoded forms) is still valid.
func expOK(v any) bool {
	if v == nil {
		return false
	}
	now := time.Now().Unix()

	switch x := v.(type) {
	case float64:
		return now <= int64(x)
	case int64:
		return now <= x
	case int:
		return now <= int64(x)
	case json.Number:
		n, err := strconv.ParseInt(string(x), 10, 64)
		if err != nil {
			return false
		}
		return now <= n
	case string:
		n, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return false
		}
		return now <= n
	default:
		return false
	}
}
