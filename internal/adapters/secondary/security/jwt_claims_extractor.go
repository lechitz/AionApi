package security

import (
	"context"
	"errors"
	"net/http"

	"github.com/lechitz/AionApi/internal/shared/common"

	"github.com/lechitz/AionApi/internal/shared/ctxkeys"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaimsExtractor extracts JWT claims from a request or context.
type JWTClaimsExtractor struct {
	secretKey string
}

// NewJWTClaimsExtractor creates and initializes a new JWTClaimsExtractor.
// secretKey is the secret key used to sign and verify JWT tokens.
func NewJWTClaimsExtractor(secretKey string) *JWTClaimsExtractor {
	return &JWTClaimsExtractor{
		secretKey: secretKey,
	}
}

// ExtractFromRequest extracts JWT claims from the request.
// Returns a map of claims or an error if the token is not found or another issue occurs.
// The token is extracted from the request cookie.
func (j *JWTClaimsExtractor) ExtractFromRequest(r *http.Request) (map[string]interface{}, error) {
	token, err := ExtractTokenFromCookie(r)
	if err != nil {
		return nil, err
	}

	claims, err := parseJWT(token, j.secretKey)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// ExtractFromContext extracts JWT claims from the context.
// Returns a map of claims or an error if the token is not found or another issue occurs.
// The token is extracted from the context.
func (j *JWTClaimsExtractor) ExtractFromContext(ctx context.Context) (map[string]interface{}, error) {
	tokenVal := ctx.Value(ctxkeys.Token)
	tokenStr, ok := tokenVal.(string)
	if !ok || tokenStr == "" {
		return nil, errors.New("token not found in context")
	}

	return parseJWT(tokenStr, j.secretKey)
}

// ExtractTokenFromCookie extracts the token from the request cookie.
// Returns the token string or an error if the token is not found or another issue occurs.
func ExtractTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(common.AuthToken)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// parseJWT parses a JWT token and returns a map of claims or an error if the token is not valid.
// secretKey is the secret key used to sign and verify JWT tokens.
func parseJWT(tokenString string, secretKey string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
