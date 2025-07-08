// Package security provides utility functions for JWT token generation and validation.
package security

import (
	"time"

	"github.com/lechitz/AionApi/internal/shared/claimskeys"

	"github.com/golang-jwt/jwt/v5"
)

// ExpTimeToken defines the duration of 24 hours used as the standard token expiration period in time-based operations.
const ExpTimeToken = 24 * time.Hour

// GenerateToken creates a signed JWT token with userID and expiration using the provided secretKey.
func GenerateToken(userID uint64, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		claimskeys.UserID: userID,
		claimskeys.Exp:    time.Now().Add(ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
