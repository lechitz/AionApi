package token_test

//
//import (
//	"errors"
//	"github.com/lechitz/AionApi/tests/setup"
//	"go.uber.org/zap/zaptest"
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	"github.com/lechitz/AionApi/internal/core/domain"
//	"github.com/lechitz/AionApi/internal/core/usecase/constants"
//	"github.com/lechitz/AionApi/internal/core/usecase/token"
//	"github.com/lechitz/AionApi/tests/mocks"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestCheck_Success(t *testing.T) {
//	suite := setup.SetupTokenServiceTest(t)
//	defer suite.Ctrl.Finish()
//
//	mockTokenRepo := mocks.NewMockTokenStore(suite.Ctrl)
//
//	logger := zaptest.NewLogger(t).Sugar()
//	secretKey := "secret"
//
//	tokenSvc := token.NewTokenService(mockTokenRepo, logger, secretKey)
//
//	ctx := domain.ContextControl{}
//	tokenString := "validToken"
//	userID := uint64(1)
//
//	suite.mockTokenRepo.EXPECT().
//		Get(ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
//		Return(tokenString, nil)
//
//	parsedUserID, parsedToken, err := tokenSvc.Check(ctx, tokenString)
//
//	assert.NoError(t, err)
//	assert.Equal(t, userID, parsedUserID)
//	assert.Equal(t, tokenString, parsedToken)
//}
//
//func TestCheck_InvalidToken(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockTokenRepo := mocks.NewMockTokenStore(ctrl)
//	logger := zaptest.NewLogger(t).Sugar()
//	secretKey := "secret"
//
//	tokenSvc := token.NewTokenService(mockTokenRepo, logger, secretKey)
//
//	ctx := domain.ContextControl{}
//	tokenString := "invalidToken"
//
//	_, _, err := tokenSvc.Check(ctx, tokenString)
//
//	assert.Error(t, err)
//	assert.Equal(t, constants.ErrorInvalidToken, err.Error())
//}
//
//func TestCheck_TokenMismatch(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockTokenRepo := mocks.NewMockTokenStore(ctrl)
//	logger := zaptest.NewLogger(t).Sugar()
//	secretKey := "secret"
//
//	tokenSvc := token.NewTokenService(mockTokenRepo, logger, secretKey)
//
//	ctx := domain.ContextControl{}
//	tokenString := "validToken"
//	userID := uint64(1)
//	cachedToken := "differentToken"
//
//	mockTokenRepo.EXPECT().
//		Get(ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
//		Return(cachedToken, nil)
//
//	_, _, err := tokenSvc.Check(ctx, tokenString)
//
//	assert.Error(t, err)
//	assert.Equal(t, constants.ErrorTokenMismatch, err.Error())
//}
//
//func TestCheck_ErrorToRetrieveTokenFromCache(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockTokenRepo := mocks.NewMockTokenStore(ctrl)
//	logger := zaptest.NewLogger(t).Sugar()
//	secretKey := "secret"
//
//	tokenSvc := token.NewTokenService(mockTokenRepo, logger, secretKey)
//
//	ctx := domain.ContextControl{}
//	tokenString := "validToken"
//	userID := uint64(1)
//
//	mockTokenRepo.EXPECT().
//		Get(ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
//		Return("", errors.New(constants.ErrorToRetrieveTokenFromCache))
//
//	_, _, err := tokenSvc.Check(ctx, tokenString)
//
//	assert.Error(t, err)
//	assert.Equal(t, constants.ErrorToRetrieveTokenFromCache, err.Error())
//}
