// Package usecase_test contains tests for user use cases.
package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
	"github.com/lechitz/AionApi/internal/user/core/ports/output"
	"github.com/lechitz/AionApi/internal/user/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := input.CreateUserCommand{
		Name:     "  User  ",
		Username: " username ",
		Email:    "  USER@example.com ",
		Password: "password",
	}

	expected := domain.User{
		ID:       1,
		Name:     "User",
		Username: "username",
		Email:    "user@example.com",
		Password: "hashed123",
		Roles:    "user",
	}

	suite.UserRepository.EXPECT().
		CheckUniqueness(gomock.Any(), "username", "user@example.com").
		Return(output.UserUniqueness{UsernameTaken: false, EmailTaken: false}, nil)

	suite.Hasher.EXPECT().
		Hash("password").
		Return("hashed123", nil)

	suite.UserRepository.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
		DoAndReturn(func(_ context.Context, u domain.User) (domain.User, error) {
			require.Equal(t, "User", u.Name)
			require.Equal(t, "username", u.Username)
			require.Equal(t, "user@example.com", u.Email)
			require.Equal(t, "hashed123", u.Password)
			u.ID = 1
			return u, nil
		})

	// (Optional) If test setup doesn't already relax logger calls:
	// suite.Logger.EXPECT().InfowCtx(gomock.Any(), usecase.SuccessUserCreated, commonkeys.UserID, gomock.Any()).AnyTimes()

	got, err := suite.UserService.Create(context.Background(), cmd)

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := input.CreateUserCommand{
		Name:     "User",
		Username: "taken",
		Email:    "user@example.com",
		Password: "password",
	}

	suite.UserRepository.EXPECT().
		CheckUniqueness(gomock.Any(), "taken", "user@example.com").
		Return(output.UserUniqueness{UsernameTaken: true, EmailTaken: false}, nil)

	got, err := suite.UserService.Create(context.Background(), cmd)

	require.Error(t, err)
	require.Equal(t, domain.User{}, got)
	require.ErrorContains(t, err, "username")
	require.ErrorContains(t, err, sharederrors.ErrUsernameInUse)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := input.CreateUserCommand{
		Name:     "User",
		Username: "user1",
		Email:    "inuse@example.com",
		Password: "password",
	}

	suite.UserRepository.EXPECT().
		CheckUniqueness(gomock.Any(), "user1", "inuse@example.com").
		Return(output.UserUniqueness{UsernameTaken: false, EmailTaken: true}, nil)

	got, err := suite.UserService.Create(context.Background(), cmd)

	require.Error(t, err)
	require.Equal(t, domain.User{}, got)
	require.ErrorContains(t, err, "email")
	require.ErrorContains(t, err, sharederrors.ErrEmailInUse)
}

func TestCreateUser_DBErrorOnCheckUniqueness(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := input.CreateUserCommand{
		Name:     "User",
		Username: "user1",
		Email:    "user1@example.com",
		Password: "password",
	}

	dbErr := errors.New("db error")
	suite.UserRepository.EXPECT().
		CheckUniqueness(gomock.Any(), "user1", "user1@example.com").
		Return(output.UserUniqueness{}, dbErr)

	got, err := suite.UserService.Create(context.Background(), cmd)

	require.Error(t, err)
	require.Equal(t, dbErr, err)
	require.Equal(t, domain.User{}, got)
}

func TestCreateUser_ErrorToHashPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := input.CreateUserCommand{
		Name:     "User",
		Username: "user1",
		Email:    "user1@example.com",
		Password: "password",
	}

	hashErr := errors.New(usecase.ErrorToHashPassword)

	suite.UserRepository.EXPECT().
		CheckUniqueness(gomock.Any(), "user1", "user1@example.com").
		Return(output.UserUniqueness{UsernameTaken: false, EmailTaken: false}, nil)

	suite.Hasher.EXPECT().
		Hash("password").
		Return("", hashErr)

	// (Optional) If the setup doesn't relax logger calls:
	// suite.Logger.EXPECT().ErrorwCtx(gomock.Any(), usecase.ErrorToHashPassword, commonkeys.Error, gomock.Any()).AnyTimes()

	got, err := suite.UserService.Create(context.Background(), cmd)

	require.Error(t, err)
	require.Equal(t, hashErr, err)
	require.Equal(t, domain.User{}, got)
}

func TestCreateUser_ErrorToCreateUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := input.CreateUserCommand{
		Name:     "User",
		Username: "user1",
		Email:    "user1@example.com",
		Password: "password",
	}

	suite.UserRepository.EXPECT().
		CheckUniqueness(gomock.Any(), "user1", "user1@example.com").
		Return(output.UserUniqueness{UsernameTaken: false, EmailTaken: false}, nil)

	suite.Hasher.EXPECT().
		Hash("password").
		Return("hashed123", nil)

	createErr := errors.New(usecase.ErrorToCreateUser)
	suite.UserRepository.EXPECT().
		Create(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
		Return(domain.User{}, createErr)

	// (Optional) If needed:
	// suite.Logger.EXPECT().ErrorwCtx(gomock.Any(), usecase.ErrorToCreateUser, commonkeys.Error, gomock.Any()).AnyTimes()

	got, err := suite.UserService.Create(context.Background(), cmd)

	require.Error(t, err)
	require.Equal(t, createErr, err)
	require.Equal(t, domain.User{}, got)
}
