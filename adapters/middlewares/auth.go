package middlewares

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			err := ValidateToken(r)
			if err != nil {
				logger.Warnw("Unauthorized access attempt", "error", err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateToken(r *http.Request) error {
	tokenString, err := extractTokenFromCookie(r)
	if err != nil {
		return err
	}
	token, err := jwt.Parse(tokenString, utils.ReturnKeyVerification)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func ExtractUserIDFromToken(r *http.Request) (uint64, error) {
	tokenString, err := extractTokenFromCookie(r)
	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenString, utils.ReturnKeyVerification)
	if err != nil {
		return 0, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", permissions["id"]), 10, 64)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

func extractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func extractTokenFromBearer(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		return strings.TrimPrefix(token, "Bearer "), nil
	}
	return "", nil
}
