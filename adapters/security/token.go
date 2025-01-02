package security

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lechitz/AionApi/config/environments"
	"time"
)

func CreateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 24).Unix()
	permissions["userId"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString(environments.Setting.SecretKey)
}

func ReturnKeyVerification(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("method signing invalid: %v", token.Header["alg"])
	}
	return environments.Setting.SecretKey, nil
}
