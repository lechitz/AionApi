// Package usecase_test contains tests for admin role management use cases.
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

func TestPromoteToAdmin_Success(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.PromoteToAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote to admin
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser},
	}

	expectedUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser, admindomain.RoleAdmin},
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, gomock.Any()).
		DoAndReturn(func(_ any, _ uint64, roles []string) (admindomain.AdminUser, error) {
			// Verify admin role was added
			require.Contains(t, roles, admindomain.RoleAdmin)
			require.Contains(t, roles, admindomain.RoleUser)
			return expectedUser, nil
		})

	result, err := suite.AdminService.PromoteToAdmin(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, expectedUser.ID, result.ID)
	require.Contains(t, result.Roles, admindomain.RoleAdmin)
}

func TestPromoteToAdmin_Unauthorized(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(2)

	cmd := admininput.PromoteToAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleUser}, // Only user role, cannot promote
	}

	// Should fail before hitting repository
	result, err := suite.AdminService.PromoteToAdmin(suite.Ctx, cmd)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrUnauthorizedPromoteToAdmin)
	require.Equal(t, uint64(0), result.ID)
}

func TestPromoteToAdmin_UserNotFound(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(999)
	actorUserID := uint64(1)

	cmd := admininput.PromoteToAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(admindomain.AdminUser{}, errors.New("user not found"))

	result, err := suite.AdminService.PromoteToAdmin(suite.Ctx, cmd)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrGetUser)
	require.Equal(t, uint64(0), result.ID)
}

func TestPromoteToAdmin_AlreadyAdmin(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.PromoteToAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser, admindomain.RoleAdmin}, // Already admin
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, currentUser.Roles).
		Return(currentUser, nil)

	result, err := suite.AdminService.PromoteToAdmin(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, targetUserID, result.ID)
	require.Contains(t, result.Roles, admindomain.RoleAdmin)
}

func TestPromoteToAdmin_RepositoryError(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.PromoteToAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser},
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, gomock.Any()).
		Return(admindomain.AdminUser{}, errors.New("database error"))

	result, err := suite.AdminService.PromoteToAdmin(suite.Ctx, cmd)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrPromoteToAdminFailed)
	require.Equal(t, uint64(0), result.ID)
}

func TestDemoteFromAdmin_Success(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.DemoteFromAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Owner can demote
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser, admindomain.RoleAdmin},
	}

	expectedUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser},
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, gomock.Any()).
		DoAndReturn(func(_ any, _ uint64, roles []string) (admindomain.AdminUser, error) {
			// Verify admin role was removed
			require.NotContains(t, roles, admindomain.RoleAdmin)
			require.Contains(t, roles, admindomain.RoleUser)
			return expectedUser, nil
		})

	result, err := suite.AdminService.DemoteFromAdmin(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, expectedUser.ID, result.ID)
	require.NotContains(t, result.Roles, admindomain.RoleAdmin)
}

func TestDemoteFromAdmin_NotAnAdmin(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.DemoteFromAdminCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner},
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser}, // Not an admin
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	// Should still succeed (idempotent)
	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, currentUser.Roles).
		Return(currentUser, nil)

	result, err := suite.AdminService.DemoteFromAdmin(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, targetUserID, result.ID)
}

func TestBlockUser_Success(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.BlockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser},
	}

	expectedUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleBlocked},
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, []string{admindomain.RoleBlocked}).
		Return(expectedUser, nil)

	result, err := suite.AdminService.BlockUser(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, expectedUser.ID, result.ID)
	require.Contains(t, result.Roles, admindomain.RoleBlocked)
}

func TestBlockUser_CannotBlockAdmin(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.BlockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleAdmin}, // Admin trying to block another admin
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser, admindomain.RoleAdmin}, // Is admin
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	result, err := suite.AdminService.BlockUser(suite.Ctx, cmd)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrUnauthorizedBlockUser)
	require.Equal(t, uint64(0), result.ID)
}

func TestBlockUser_AlreadyBlocked(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.BlockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleBlocked}, // Already blocked
	}

	suite.AdminRepository.EXPECT().
		GetByID(gomock.Any(), targetUserID).
		Return(currentUser, nil)

	// Should still update (idempotent)
	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, []string{admindomain.RoleBlocked}).
		Return(currentUser, nil)

	result, err := suite.AdminService.BlockUser(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, targetUserID, result.ID)
	require.Contains(t, result.Roles, admindomain.RoleBlocked)
}

func TestUnblockUser_Success(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.UnblockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	expectedUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser},
	}

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, []string{admindomain.RoleUser}).
		Return(expectedUser, nil)

	result, err := suite.AdminService.UnblockUser(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, expectedUser.ID, result.ID)
	require.NotContains(t, result.Roles, admindomain.RoleBlocked)
	require.Contains(t, result.Roles, admindomain.RoleUser)
}

func TestUnblockUser_NotBlocked(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.UnblockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	currentUser := admindomain.AdminUser{
		ID:    targetUserID,
		Roles: []string{admindomain.RoleUser}, // Not blocked
	}

	// Should still succeed (idempotent)
	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, []string{admindomain.RoleUser}).
		Return(currentUser, nil)

	result, err := suite.AdminService.UnblockUser(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, targetUserID, result.ID)
}

func TestUnblockUser_RepositoryError(t *testing.T) {
	suite := setup.AdminServiceTest(t)
	defer suite.Ctrl.Finish()

	targetUserID := uint64(100)
	actorUserID := uint64(1)

	cmd := admininput.UnblockUserCommand{
		UserID:      targetUserID,
		ActorUserID: actorUserID,
		ActorRoles:  []string{admindomain.RoleOwner}, // Must be owner to promote/demote admin
	}

	suite.AdminRepository.EXPECT().
		UpdateRoles(gomock.Any(), targetUserID, gomock.Any()).
		Return(admindomain.AdminUser{}, errors.New("database error"))

	result, err := suite.AdminService.UnblockUser(suite.Ctx, cmd)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrUnblockUserFailed)
	require.Equal(t, uint64(0), result.ID)
}
