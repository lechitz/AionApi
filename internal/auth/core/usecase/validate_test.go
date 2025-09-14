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

// --- Happy paths ---

// Valid token with canonical userID ("userID") as string → success.
func TestValidate_Success_UserIDString(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 123
	const raw = "Bearer abc.def.ghi"

	suite.AuthProvider.EXPECT().
		Verify("abc.def.ghi").
		Return(map[string]any{claimskeys.UserID: "123"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID).
		Return(domain.Auth{Key: userID, Token: "abc.def.ghi"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userID, uid)
	require.Equal(t, "123", claims[claimskeys.UserID])
}

// Claims use "sub" instead of "userID" with spaces → success, also exercises TrimSpace.
func TestValidate_Success_SubWithSpaces(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.sub.9"

	suite.AuthProvider.EXPECT().
		Verify("tok.sub.9").
		Return(map[string]any{"sub": "   9   "}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(9)).
		Return(domain.Auth{Key: 9, Token: "tok.sub.9"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, uint64(9), uid)
	require.Equal(t, "   9   ", claims["sub"])
}

// userID as json.Number → success.
func TestValidate_Success_JSONNumberClaim(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.json.123"

	suite.AuthProvider.EXPECT().
		Verify("tok.json.123").
		Return(map[string]any{claimskeys.UserID: json.Number("123")}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(123)).
		Return(domain.Auth{Key: 123, Token: "tok.json.123"}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, uint64(123), uid)
}

// userID como float64 integral → sucesso.
func TestValidate_Success_Float64Integral(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.f64.777"

	suite.AuthProvider.EXPECT().
		Verify("tok.f64.777").
		Return(map[string]any{claimskeys.UserID: float64(777)}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(777)).
		Return(domain.Auth{Key: 777, Token: "tok.f64.777"}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, uint64(777), uid)
}

// Sanitização: prefixo Bearer “case-insensitive” + espaços → sucesso.
func TestValidate_Success_SanitizeBearerAndSpaces(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "   bEaReR tok.sanitized   "

	suite.AuthProvider.EXPECT().
		Verify("tok.sanitized").
		Return(map[string]any{claimskeys.UserID: "88"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(88)).
		Return(domain.Auth{Key: 88, Token: "tok.sanitized"}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, uint64(88), uid)
}

// Sanitização: sem “Bearer”, apenas espaços → sucesso.
func TestValidate_Success_SanitizeOnlySpaces_NoBearer(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "   tok.no.bearer   "

	suite.AuthProvider.EXPECT().
		Verify("tok.no.bearer").
		Return(map[string]any{claimskeys.UserID: "11"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(11)).
		Return(domain.Auth{Key: 11, Token: "tok.no.bearer"}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, uint64(11), uid)
}

// --- Error paths ---

// Invalid token (provider verify falha) → ErrorInvalidToken.
func TestValidate_Error_InvalidToken(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer bad.token"

	suite.AuthProvider.EXPECT().
		Verify("bad.token").
		Return(nil, errors.New(usecase.ErrorInvalidToken))

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Equal(t, usecase.ErrorInvalidToken, err.Error())
}

// userID ausente (nem userID nem sub) → ErrorInvalidUserIDClaim.
func TestValidate_Error_MissingUserIDAndSub(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.missing.claims"

	suite.AuthProvider.EXPECT().
		Verify("tok.missing.claims").
		Return(map[string]any{"foo": "bar"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Equal(t, usecase.ErrorInvalidUserIDClaim, err.Error())
}

// userID como float64 **não integral** → ErrorInvalidUserIDClaim.
func TestValidate_Error_Float64NonIntegral(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.f64.bad"

	suite.AuthProvider.EXPECT().
		Verify("tok.f64.bad").
		Return(map[string]any{claimskeys.UserID: 123.5}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Equal(t, usecase.ErrorInvalidUserIDClaim, err.Error())
}

// userID como float64 **negativo** → ErrorInvalidUserIDClaim (pega ramo t<0).
func TestValidate_Error_Float64Negative(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.f64.neg"

	suite.AuthProvider.EXPECT().
		Verify("tok.f64.neg").
		Return(map[string]any{claimskeys.UserID: -1.0}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Equal(t, usecase.ErrorInvalidUserIDClaim, err.Error())
}

// userID como tipo **não suportado** (int) → ErrorInvalidUserIDClaim (cai no default do switch).
func TestValidate_Error_UserIDUnsupportedType(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.unsupported.type"

	suite.AuthProvider.EXPECT().
		Verify("tok.unsupported.type").
		Return(map[string]any{claimskeys.UserID: int(123)}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Equal(t, usecase.ErrorInvalidUserIDClaim, err.Error())
}

// Store.Get retorna erro → ErrorToRetrieveTokenFromCache.
func TestValidate_Error_StoreGet(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.store.err"

	suite.AuthProvider.EXPECT().
		Verify("tok.store.err").
		Return(map[string]any{claimskeys.UserID: "42"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(42)).
		Return(domain.Auth{}, errors.New(usecase.ErrorToRetrieveTokenFromCache))

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Equal(t, usecase.ErrorToRetrieveTokenFromCache, err.Error())
}

// Token/Cache mismatch → ErrorTokenMismatch.
func TestValidate_Error_Mismatch(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 5
	const raw = "Bearer abc.def.ghi"

	suite.AuthProvider.EXPECT().
		Verify("abc.def.ghi").
		Return(map[string]any{claimskeys.UserID: "5"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID).
		Return(domain.Auth{Key: userID, Token: "zzz.yyy.xxx"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.ErrorContains(t, err, usecase.ErrorTokenMismatch)
}

// Cache vazio (token vazio) → ErrorTokenMismatch.
func TestValidate_Error_EmptyTokenInCache(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.empty.cache"

	suite.AuthProvider.EXPECT().
		Verify("tok.empty.cache").
		Return(map[string]any{claimskeys.UserID: "77"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(77)).
		Return(domain.Auth{Key: 77, Token: ""}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Equal(t, usecase.ErrorTokenMismatch, err.Error())
}

func TestValidate_Error_SubClaimInvalidJSONNumber(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer tok.sub.bad"

	// Força uso do ramo "sub" (não userID) com um valor inválido
	suite.AuthProvider.EXPECT().
		Verify("tok.sub.bad").
		Return(map[string]any{"sub": json.Number("abc")}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, usecase.ErrorInvalidUserIDClaim, err.Error())
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
}

// Exercita o sanitizeTokenValue no caminho do "return s":
// token já limpo (sem espaços) e sem prefixo "Bearer ".
// Isso garante cobertura da linha final do sanitizador.
func TestValidate_Success_NoBearerNoSpaces_ReturnS(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "tok.clean"

	suite.AuthProvider.EXPECT().
		Verify("tok.clean").
		Return(map[string]any{claimskeys.UserID: "22"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), uint64(22)).
		Return(domain.Auth{Key: 22, Token: "tok.clean"}, nil)

	uid, _, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, uint64(22), uid)
}
