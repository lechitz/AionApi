// Package jwt provides utility functions for JWT token generation and validation.
package jwt

import (
	"github.com/lechitz/AionApi/internal/adapters/secondary/security/jwt/constants"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT token with userID and expiration using the provided secretKey.
func GenerateToken(userID uint64, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		claimskeys.UserID: userID,
		claimskeys.Exp:    time.Now().Add(constants.ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func ParseToken() {}
