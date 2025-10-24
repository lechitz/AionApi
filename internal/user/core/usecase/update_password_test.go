// File: internal/user/core/usecase/update_password_test.go
package usecase_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	authdomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/internal/user/core/usecase"

	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newUserService(t *testing.T) (context.Context,
	*mocks.MockUserRepository,
	*mocks.MockHasher,
	*mocks.MockAuthProvider,
	*mocks.MockAuthStore,
	*usecase.Service,
	*gomock.Controller,
) {
	t.Helper()

	ctrl := gomock.NewController(t)
	ctx := t.Context()

	repo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	provider := mocks.NewMockAuthProvider(ctrl)
	store := mocks.NewMockAuthStore(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)

	// Relaxed expectations for logger to tolerate different arities.
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()                             // 4 args (e.g., empty token path)
	logger.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes() // 6 args

	svc := usecase.NewService(
		repo,  // userRepository
		store, // authStore
		provider,
		hasher,
		logger,
	)

	return ctx, repo, hasher, provider, store, svc, ctrl
}

func defaultUser() userdomain.User {
	now := time.Now().UTC()
	return userdomain.User{
		ID:        42,
		Name:      "Test User",
		Username:  "testuser",
		Email:     "user@example.com",
		Password:  "stored-hash",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestUpdatePassword_Success(t *testing.T) {
	ctx, repo, hasher, provider, store, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	newHash := "new-hash"
	token := "new-token"

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	hasher.EXPECT().Compare(u.Password, "oldPassword").Return(nil)
	hasher.EXPECT().Hash("newPassword").Return(newHash, nil)
	repo.EXPECT().
		Update(gomock.Any(), u.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(u, nil)
	provider.EXPECT().GenerateRefreshToken(u.ID).Return(token, nil)

	// The implementation calls Verify to compute TTL before saving; expect it and return a positive exp
	exp := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)
	provider.EXPECT().Verify(token).Return(map[string]any{claimskeys.Exp: exp}, nil)

	store.EXPECT().Save(gomock.Any(), authdomain.Auth{Key: u.ID, Token: token, Type: commonkeys.TokenTypeAccess}, gomock.Any()).Return(nil)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.NoError(t, err)
	require.Equal(t, token, gotToken)
}

func TestUpdatePassword_ErrorToGetSelf(t *testing.T) {
	ctx, repo, _, _, _, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	expected := errors.New(usecase.ErrorToGetSelf)

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(userdomain.User{}, expected)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.Error(t, err)
	require.Empty(t, gotToken)
	require.Contains(t, err.Error(), usecase.ErrorToGetSelf)
}

func TestUpdatePassword_ErrorToCompareHashAndPassword(t *testing.T) {
	ctx, repo, hasher, _, _, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	expected := errors.New(usecase.ErrorToCompareHashAndPassword)

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	hasher.EXPECT().Compare(u.Password, "oldPassword").Return(expected)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.Error(t, err)
	require.Empty(t, gotToken)
	require.Contains(t, err.Error(), usecase.ErrorToCompareHashAndPassword)
}

func TestUpdatePassword_ErrorToHashPassword(t *testing.T) {
	ctx, repo, hasher, _, _, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	expected := errors.New(usecase.ErrorToHashPassword)

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	hasher.EXPECT().Compare(u.Password, "oldPassword").Return(nil)
	hasher.EXPECT().Hash("newPassword").Return("", expected)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.Error(t, err)
	require.Empty(t, gotToken)
	require.Contains(t, err.Error(), usecase.ErrorToHashPassword)
}

func TestUpdatePassword_ErrorToUpdatePassword(t *testing.T) {
	ctx, repo, hasher, _, _, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	newHash := "new-hash"
	expected := errors.New(usecase.ErrorToUpdatePassword)

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	hasher.EXPECT().Compare(u.Password, "oldPassword").Return(nil)
	hasher.EXPECT().Hash("newPassword").Return(newHash, nil)
	repo.EXPECT().
		Update(gomock.Any(), u.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(userdomain.User{}, expected)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.Error(t, err)
	require.Empty(t, gotToken)
	require.Contains(t, err.Error(), usecase.ErrorToUpdatePassword)
}

func TestUpdatePassword_ErrorToCreateToken_WithProviderError(t *testing.T) {
	ctx, repo, hasher, provider, _, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	newHash := "new-hash"
	providerErr := errors.New(usecase.ErrorToCreateToken)

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	hasher.EXPECT().Compare(u.Password, "oldPassword").Return(nil)
	hasher.EXPECT().Hash("newPassword").Return(newHash, nil)
	repo.EXPECT().
		Update(gomock.Any(), u.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(u, nil)
	provider.EXPECT().GenerateRefreshToken(u.ID).Return("", providerErr)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.Error(t, err)
	require.Empty(t, gotToken)
	require.Contains(t, err.Error(), usecase.ErrorToCreateToken)
}

func TestUpdatePassword_ErrorToCreateToken_EmptyToken(t *testing.T) {
	ctx, repo, hasher, provider, _, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	newHash := "new-hash"

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	hasher.EXPECT().Compare(u.Password, "oldPassword").Return(nil)
	hasher.EXPECT().Hash("newPassword").Return(newHash, nil)
	repo.EXPECT().
		Update(gomock.Any(), u.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(u, nil)
	// Provider returns empty token without error.
	provider.EXPECT().GenerateRefreshToken(u.ID).Return("", nil)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.Error(t, err)
	require.Empty(t, gotToken)
	require.Contains(t, err.Error(), usecase.ErrorToCreateToken)
}

func TestUpdatePassword_ErrorToSaveToken(t *testing.T) {
	ctx, repo, hasher, provider, store, svc, ctrl := newUserService(t)
	defer ctrl.Finish()

	u := defaultUser()
	newHash := "new-hash"
	token := "new-token"
	saveErr := errors.New(usecase.ErrorToCreateToken)

	repo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	hasher.EXPECT().Compare(u.Password, "oldPassword").Return(nil)
	hasher.EXPECT().Hash("newPassword").Return(newHash, nil)
	repo.EXPECT().
		Update(gomock.Any(), u.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(u, nil)
	provider.EXPECT().GenerateRefreshToken(u.ID).Return(token, nil)

	// Expect Verify to compute TTL; return a future exp so Save is attempted and returns the configured error
	exp := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)
	provider.EXPECT().Verify(token).Return(map[string]any{claimskeys.Exp: exp}, nil)

	store.EXPECT().Save(gomock.Any(), authdomain.Auth{Key: u.ID, Token: token, Type: commonkeys.TokenTypeAccess}, gomock.Any()).Return(saveErr)

	gotToken, err := svc.UpdatePassword(ctx, u.ID, "oldPassword", "newPassword")
	require.Error(t, err)
	require.Empty(t, gotToken)
	require.Contains(t, err.Error(), usecase.ErrorToCreateToken)
}
