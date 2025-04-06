package auth_test

//
//import (
//	"errors"
//	"github.com/lechitz/AionApi/tests/mocks"
//	"testing"
//
//	"github.com/lechitz/AionApi/internal/core/domain"
//	"github.com/lechitz/AionApi/tests/setup"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestLogout_Success(t *testing.T) {
//	suite := setup.SetupAuthServiceTest(t)
//	defer suite.Ctrl.Finish()
//
//	token := "valid.token.value"
//	userID := uint64(1)
//	tokenDomain := domain.TokenDomain{
//		UserID: userID,
//		Token:  token,
//	}
//
//	mockTokenService := mocks.NewMockTokenStore(suite.Ctrl)
//
//	mockTokenService.EXPECT().
//		Check(suite.Ctx, token).
//		Return(userID, token, nil)
//
//	mockTokenService.EXPECT().
//		Delete(suite.Ctx, tokenDomain).Return(nil)
//
//	err := suite.AuthService.Logout(suite.Ctx, token)
//
//	assert.NoError(t, err)
//}
//
//func TestLogout_CheckTokenFails(t *testing.T) {
//	suite := setup.SetupAuthServiceTest(t)
//	defer suite.Ctrl.Finish()
//
//	token := "invalid.token.value"
//	expectedErr := errors.New("invalid token")
//
//	mockTokenService := mocks.NewMockTokenStore(suite.Ctrl)
//
//	mockTokenService.EXPECT().
//		Check(suite.Ctx, token).
//		Return(uint64(0), "", expectedErr)
//
//	err := suite.AuthService.Logout(suite.Ctx, token)
//
//	assert.EqualError(t, err, expectedErr.Error())
//}
//
//func TestLogout_DeleteTokenFails(t *testing.T) {
//	suite := setup.SetupAuthServiceTest(t)
//	defer suite.Ctrl.Finish()
//
//	token := "valid.token.value"
//	userID := uint64(1)
//	expectedErr := errors.New("delete error")
//
//	mockTokenService := mocks.NewMockTokenStore(suite.Ctrl)
//
//	mockTokenService.EXPECT().
//		Check(suite.Ctx, token).
//		Return(userID, token, nil)
//
//	mockTokenService.EXPECT().
//		Delete(suite.Ctx, domain.TokenDomain{UserID: userID, Token: token}).
//		Return(expectedErr)
//
//	err := suite.AuthService.Logout(suite.Ctx, token)
//
//	assert.EqualError(t, err, expectedErr.Error())
//}
