package usecase_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/auth/core/usecase"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// -----------------------------------------------------------------------------
// Helpers & shared constants
// -----------------------------------------------------------------------------

const (
	bearerPrefix = "Bearer "
)

func requireUnauthorizedWith(t *testing.T, err error, contains string) {
	t.Helper()
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")
	require.Contains(t, err.Error(), contains)
}

func requireZeroAndNilClaims(t *testing.T, uid uint64, claims map[string]any) {
	t.Helper()
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
}

// -----------------------------------------------------------------------------
// Happy paths
// -----------------------------------------------------------------------------

// userID present as string in claims → success.
func TestValidate_Success_UserIDString(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		userID = uint64(123)
		token  = "abc.def.ghi"
		raw    = bearerPrefix + token
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: "123"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID).
		Return(domain.Auth{Key: userID, Token: token}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userID, uid)
	require.Equal(t, "123", claims[claimskeys.UserID])
}

// Fallback to "sub" with spaces (TrimSpace exercised) → success.
func TestValidate_Success_SubWithSpaces(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token      = "tok.sub.9"
		raw        = bearerPrefix + token
		subRaw     = "   9   "
		userIDWant = uint64(9)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{"sub": subRaw}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: token}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userIDWant, uid)
	require.Equal(t, subRaw, claims["sub"])
}

// userID as json.Number → success.
func TestValidate_Success_JSONNumberClaim(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token      = "tok.json.123"
		raw        = bearerPrefix + token
		userIDWant = uint64(123)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: json.Number("123")}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: token}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userIDWant, uid)
}

// userID as integral float64 → success.
func TestValidate_Success_Float64Integral(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token      = "tok.f64.777"
		raw        = bearerPrefix + token
		userIDWant = uint64(777)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: float64(777)}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: token}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userIDWant, uid)
}

// Sanitization: case-insensitive "Bearer" + surrounding spaces → success.
func TestValidate_Success_SanitizeBearerAndSpaces(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		rawMessy   = "   bEaReR tok.sanitized   "
		token      = "tok.sanitized" //nolint:gosec // fake token for tests
		userIDWant = uint64(88)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: "88"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: token}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, rawMessy)
	require.NoError(t, err)
	require.Equal(t, userIDWant, uid)
}

// Sanitization: only spaces, no "Bearer" prefix → success.
func TestValidate_Success_SanitizeOnlySpaces_NoBearer(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		rawMessy   = "   tok.no.bearer   "
		token      = "tok.no.bearer"
		userIDWant = uint64(11)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: "11"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: token}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, rawMessy)
	require.NoError(t, err)
	require.Equal(t, userIDWant, uid)
}

// Fast path: already clean, no spaces/prefix → success.
func TestValidate_Success_NoBearerNoSpaces_ReturnS(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token      = "tok.clean"
		raw        = token
		userIDWant = uint64(22)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: "22"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: token}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userIDWant, uid)
}

// -----------------------------------------------------------------------------
// Error paths
// -----------------------------------------------------------------------------

// Verify fails → unauthorized with ErrorInvalidToken in message.
func TestValidate_Error_InvalidToken(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token = "bad.token"
		raw   = bearerPrefix + token
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(nil, errors.New(usecase.ErrorInvalidToken))

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorInvalidToken)
	requireZeroAndNilClaims(t, uid, claims)
}

// Missing userID in claims (neither "userID" nor "sub") → unauthorized with ErrorInvalidUserIDClaim.
func TestValidate_Error_MissingUserIDAndSub(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token = "tok.missing.claims"
		raw   = bearerPrefix + token
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{"foo": "bar"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorInvalidUserIDClaim)
	requireZeroAndNilClaims(t, uid, claims)
}

// userID as non-integral float64 → unauthorized with ErrorInvalidUserIDClaim.
func TestValidate_Error_Float64NonIntegral(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token = "tok.f64.bad" //nolint:gosec // fake token for tests
		raw   = bearerPrefix + token
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: 123.5}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorInvalidUserIDClaim)
	requireZeroAndNilClaims(t, uid, claims)
}

// userID as negative float64 → unauthorized with ErrorInvalidUserIDClaim.
func TestValidate_Error_Float64Negative(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token = "tok.f64.neg" //nolint:gosec // fake token for tests
		raw   = bearerPrefix + token
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: -1.0}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorInvalidUserIDClaim)
	requireZeroAndNilClaims(t, uid, claims)
}

// userID as unsupported type (int) → unauthorized with ErrorInvalidUserIDClaim.
func TestValidate_Error_UserIDUnsupportedType(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token = "tok.unsupported.type"
		raw   = bearerPrefix + token
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: int(123)}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorInvalidUserIDClaim)
	requireZeroAndNilClaims(t, uid, claims)
}

// Store.Get returns error → raw ErrorToRetrieveTokenFromCache (not unauthorized).
func TestValidate_Error_StoreGet(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token      = "tok.store.err" //nolint:gosec // simulated cache error token
		raw        = bearerPrefix + token
		userIDWant = uint64(42)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: "42"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{}, errors.New(usecase.ErrorToRetrieveTokenFromCache))

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, usecase.ErrorToRetrieveTokenFromCache, err.Error())
	requireZeroAndNilClaims(t, uid, claims)
}

// Token/cache mismatch → unauthorized with ErrorTokenMismatch.
func TestValidate_Error_Mismatch(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		userIDWant           = uint64(5)
		token                = "abc.def.ghi"
		raw                  = bearerPrefix + token
		cachedTokenDifferent = "zzz.yyy.xxx"
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: "5"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: cachedTokenDifferent}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorTokenMismatch)
	requireZeroAndNilClaims(t, uid, claims)
}

// Empty token in cache → unauthorized with ErrorTokenMismatch.
func TestValidate_Error_EmptyTokenInCache(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token      = "tok.empty.cache"
		raw        = bearerPrefix + token
		userIDWant = uint64(77)
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{claimskeys.UserID: "77"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userIDWant).
		Return(domain.Auth{Key: userIDWant, Token: ""}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorTokenMismatch)
	requireZeroAndNilClaims(t, uid, claims)
}

// "sub" is present but invalid json.Number → unauthorized with ErrorInvalidUserIDClaim.
func TestValidate_Error_SubClaimInvalidJSONNumber(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const (
		token = "tok.sub.bad"
		raw   = bearerPrefix + token
	)

	suite.AuthProvider.EXPECT().
		Verify(token).
		Return(map[string]any{"sub": json.Number("abc")}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	requireUnauthorizedWith(t, err, usecase.ErrorInvalidUserIDClaim)
	requireZeroAndNilClaims(t, uid, claims)
}
