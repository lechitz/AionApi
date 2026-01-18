// File: internal/admin/core/usecase/update_test.go
package usecase_test

import (
	"errors"
	"testing"

	admindomain "github.com/lechitz/AionApi/internal/admin/core/domain"
	admininput "github.com/lechitz/AionApi/internal/admin/core/ports/input"
	"github.com/lechitz/AionApi/internal/admin/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateUserRoles_Success(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := uint64(1)
	cmd := admininput.UpdateUserRolesCommand{
		UserID: uid,
		Roles:  []string{"user", "admin"},
	}

	currentUser := setup.DefaultTestUserWithRoles([]string{"user"})
	expectedUser := setup.DefaultTestUserWithRoles([]string{"user", "admin"})

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), uid).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), uid, cmd.Roles).
		Return(expectedUser, nil)

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, expectedUser.ID, got.ID)
	require.Equal(t, expectedUser.Roles, got.Roles)
}

func TestUpdateUserRoles_BlockUser(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := uint64(2)
	cmd := admininput.UpdateUserRolesCommand{
		UserID: uid,
		Roles:  []string{"blocked"},
	}

	currentUser := setup.DefaultTestUserWithRoles([]string{"user"})
	expectedUser := setup.DefaultTestUserWithRoles([]string{"blocked"})
	expectedUser.ID = uid

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), uid).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), uid, cmd.Roles).
		Return(expectedUser, nil)

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, []string{"blocked"}, got.Roles)
}

func TestUpdateUserRoles_NoRolesToUpdate(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := admininput.UpdateUserRolesCommand{
		UserID: 1,
		Roles:  []string{}, // empty roles
	}

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.Error(t, err)
	require.Equal(t, admindomain.AdminUser{}, got)
	require.ErrorIs(t, err, usecase.ErrNoRolesToUpdate)
}

func TestUpdateUserRoles_InvalidRole(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	cmd := admininput.UpdateUserRolesCommand{
		UserID: 1,
		Roles:  []string{"invalid_role"},
	}

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.Error(t, err)
	require.Equal(t, admindomain.AdminUser{}, got)
	require.ErrorContains(t, err, usecase.ErrorInvalidRole)
}

func TestUpdateUserRoles_CannotBlockAdmin(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := uint64(1)
	cmd := admininput.UpdateUserRolesCommand{
		UserID: uid,
		Roles:  []string{"blocked"},
	}

	// User is an admin
	currentUser := setup.DefaultTestUserWithRoles([]string{"admin"})

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), uid).
		Return(currentUser, nil)

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.Error(t, err)
	require.Equal(t, admindomain.AdminUser{}, got)
	require.ErrorIs(t, err, usecase.ErrCannotBlockAdmin)
}

func TestUpdateUserRoles_UserNotFound(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := uint64(999)
	cmd := admininput.UpdateUserRolesCommand{
		UserID: uid,
		Roles:  []string{"user"},
	}

	repoErr := errors.New("record not found")

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), uid).
		Return(admindomain.AdminUser{}, repoErr)

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.Error(t, err)
	require.Equal(t, admindomain.AdminUser{}, got)
	require.ErrorContains(t, err, usecase.ErrorToGetUser)
}

func TestUpdateUserRoles_RepositoryError(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := uint64(1)
	cmd := admininput.UpdateUserRolesCommand{
		UserID: uid,
		Roles:  []string{"user", "admin"},
	}

	currentUser := setup.DefaultTestUserWithRoles([]string{"user"})
	repoErr := errors.New("database error")

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), uid).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), uid, cmd.Roles).
		Return(admindomain.AdminUser{}, repoErr)

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.Error(t, err)
	require.Equal(t, admindomain.AdminUser{}, got)
	require.ErrorContains(t, err, usecase.ErrorToUpdateUserRoles)
}

func TestUpdateUserRoles_PromoteToAdmin(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := uint64(3)
	cmd := admininput.UpdateUserRolesCommand{
		UserID: uid,
		Roles:  []string{"admin"},
	}

	currentUser := setup.DefaultTestUserWithRoles([]string{"user"})
	currentUser.ID = uid
	expectedUser := setup.DefaultTestUserWithRoles([]string{"admin"})
	expectedUser.ID = uid

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), uid).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), uid, cmd.Roles).
		Return(expectedUser, nil)

	got, err := suite.AdminService.UpdateUserRoles(suite.Ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, []string{"admin"}, got.Roles)
}
