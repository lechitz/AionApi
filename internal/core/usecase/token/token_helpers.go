package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/platform/config"
	"time"
)

func generateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		constants.UserID: userID,
		constants.Exp:    time.Now().Add(constants.ExpTimeToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Setting.SecretKey))
}
