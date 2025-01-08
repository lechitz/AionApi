package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/config"
	"time"
)

func CreateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
		"id":         userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString(config.Setting.SecretKey)
}

func ReturnKeyVerification(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("method signing invalid: %v", token.Header["alg"])
	}
	return config.Setting.SecretKey, nil
}
