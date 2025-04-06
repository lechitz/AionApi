package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

func GenerateToken(userID uint64, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		constants.UserID: userID,
		constants.Exp:    time.Now().Add(constants.ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
