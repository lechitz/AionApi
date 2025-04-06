package user_test

//import (
//	"errors"
//	"github.com/lechitz/AionApi/internal/core/domain"
//	"github.com/lechitz/AionApi/internal/core/usecase/constants"
//	"github.com/lechitz/AionApi/internal/core/usecase/user"
//	"github.com/lechitz/AionApi/tests/setup"
//	"github.com/stretchr/testify/assert"
//	"go.uber.org/zap/zaptest"
//	"testing"
//)
//
//func TestDeleteUser_Success(t *testing.T) {
//	suite := setup.SetupUserServiceTest(t)
//	defer suite.Ctrl.Finish()
//
//	ctx := domain.ContextControl{}
//	userID := uint64(1)
//
//	suite.UserStore.EXPECT().
//		SoftDeleteUser(ctx, userID).
//		Return(nil)
//
//	suite.TokenSvc.EXPECT().
//		Delete(ctx, domain.TokenDomain{UserID: userID}).
//		Return(nil)
//
//	err := userSvc.SoftDeleteUser(ctx, userID)
//
//	assert.NoError(t, err)
//}
//
//func TestDeleteUser_ErrorToSoftDeleteUser(t *testing.T) {
//	suite := setup.SetupUserServiceTest(t)
//	defer suite.Ctrl.Finish()
//
//	logger := zaptest.NewLogger(t).Sugar()
//
//	userSvc := user.NewUserService(suite.UserStore, suite.TokenService, suite.HasherStore, logger)
//
//	ctx := domain.ContextControl{}
//	userID := uint64(1)
//	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)
//
//	suite.UserStore.EXPECT().
//		SoftDeleteUser(ctx, userID).
//		Return(expectedErr)
//
//	err := userSvc.SoftDeleteUser(ctx, userID)
//
//	assert.EqualError(t, err, expectedErr.Error())
//}
//
//func TestDeleteUser_ErrorToDeleteToken(t *testing.T) {
//	suite := setup.SetupUserServiceTest(t)
//	defer suite.Ctrl.Finish()
//
//	logger := zaptest.NewLogger(t).Sugar()
//
//	userSvc := user.NewUserService(suite.UserStore, suite.TokenSvc, suite.HasherStore, logger)
//
//	ctx := domain.ContextControl{}
//	userID := uint64(1)
//	expectedErr := errors.New(constants.ErrorToDeleteToken)
//
//	suite.UserStore.EXPECT().
//		SoftDeleteUser(ctx, userID).
//		Return(nil)
//
//	suite.TokenSvc.EXPECT().
//		Delete(ctx, domain.TokenDomain{UserID: userID}).
//		Return(expectedErr)
//
//	err := userSvc.SoftDeleteUser(ctx, userID)
//
//	assert.EqualError(t, err, expectedErr.Error())
//}
