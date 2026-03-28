package domain_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/admin/core/domain"
	"github.com/stretchr/testify/require"
)

func TestAdminUserRoleHelpers(t *testing.T) {
	user := domain.AdminUser{Roles: []string{domain.RoleUser, domain.RoleAdmin}}

	require.True(t, user.HasRole(domain.RoleAdmin))
	require.Equal(t, domain.RoleAdmin, user.GetHighestRole())
	require.False(t, user.IsBlocked())
	require.True(t, user.IsAdmin())
	require.False(t, user.IsOwner())
}

func TestValidRolesAndIsValidRole(t *testing.T) {
	require.Equal(t, []string{domain.RoleOwner, domain.RoleAdmin, domain.RoleUser, domain.RoleBlocked}, domain.ValidRoles())
	require.True(t, domain.IsValidRole(domain.RoleOwner))
	require.False(t, domain.IsValidRole("invalid"))
}

func TestRoleHierarchyFunctions(t *testing.T) {
	require.True(t, domain.CanManageRole(domain.RoleOwner, domain.RoleAdmin))
	require.False(t, domain.CanManageRole(domain.RoleAdmin, domain.RoleOwner))
	require.False(t, domain.CanManageRole("unknown", domain.RoleUser))
	require.True(t, domain.HasRole([]string{domain.RoleUser, domain.RoleBlocked}, domain.RoleBlocked))
	require.False(t, domain.HasRole([]string{domain.RoleUser}, domain.RoleAdmin))
	require.Equal(t, domain.RoleBlocked, domain.GetHighestRole([]string{"unknown"}))
	require.Equal(t, domain.RoleOwner, domain.GetHighestRole([]string{domain.RoleUser, domain.RoleOwner, domain.RoleAdmin}))
}
