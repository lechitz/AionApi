package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	UserID       = "user_id"
	Exp          = "exp"
	ExpTimeToken = 24 * time.Hour
)

func GenerateToken(userID uint64, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		UserID: userID,
		Exp:    time.Now().Add(ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
