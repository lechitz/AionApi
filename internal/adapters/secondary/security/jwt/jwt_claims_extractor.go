package jwt

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"net/http"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"

	"github.com/golang-jwt/jwt/v5"
)

// ClaimsExtractor extracts JWT claims from a request or context.
type ClaimsExtractor struct {
	secretKey string
}

// NewClaimsExtractor creates and initializes a new ClaimsExtractor.
// secretKey is the secret key used to sign and verify JWT tokens.
func NewClaimsExtractor(secretKey string) *ClaimsExtractor {
	return &ClaimsExtractor{
		secretKey: secretKey,
	}
}

// ExtractFromRequest extracts JWT claims from the request.
// Returns a map of claims or an error if the token is not found or another issue occurs.
// The token is extracted from the request cookie.
func (j *ClaimsExtractor) ExtractFromRequest(r *http.Request) (map[string]interface{}, error) {
	token, err := extractTokenFromCookie(r)
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
func (j *ClaimsExtractor) ExtractFromContext(ctx context.Context) (map[string]interface{}, error) {
	tokenVal := ctx.Value(ctxkeys.Token)
	tokenStr, ok := tokenVal.(string)
	if !ok || tokenStr == "" {
		return nil, errors.New(sharederrors.ErrTokenNotFound)
	}

	return parseJWT(tokenStr, j.secretKey)
}

// ExtractTokenFromCookie extracts the token from the request cookie.
// Returns the token string or an error if the token is not found or another issue occurs.
func extractTokenFromCookie(r *http.Request) (string, error) { //TODO: está sendo usada em auth_middleware atualmente.
	cookie, err := r.Cookie(commonkeys.AuthTokenCookieName)
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
		return nil, errors.New(sharederrors.ErrInvalidToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(sharederrors.ErrInvalidClaims)
	}

	return claims, nil
}
