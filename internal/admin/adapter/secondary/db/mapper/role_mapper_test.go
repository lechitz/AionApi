package mapper_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/admin/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/admin/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/admin/core/domain"
	"github.com/stretchr/testify/require"
)

func TestRoleRoundtrip(t *testing.T) {
	now := time.Now().UTC()
	desc := "administrator"

	dbRole := model.RoleDB{ID: 1, Name: domain.RoleAdmin, Description: &desc, IsActive: true, CreatedAt: now, UpdatedAt: now}
	dRole := mapper.RoleFromDB(dbRole)
	require.Equal(t, desc, dRole.Description)

	backToDB := mapper.RoleToDB(dRole)
	require.NotNil(t, backToDB.Description)
	require.Equal(t, desc, *backToDB.Description)

	dRole.Description = ""
	withoutDesc := mapper.RoleToDB(dRole)
	require.Nil(t, withoutDesc.Description)
}

func TestUserRoleRoundtrip(t *testing.T) {
	now := time.Now().UTC()
	assignedBy := uint64(123)
	dbUserRole := model.UserRoleDB{ID: 4, UserID: 8, RoleID: 2, AssignedBy: &assignedBy, AssignedAt: now}

	domainUserRole := mapper.UserRoleFromDB(dbUserRole)
	require.Equal(t, dbUserRole.ID, domainUserRole.ID)
	require.Equal(t, dbUserRole.AssignedBy, domainUserRole.AssignedBy)

	backToDB := mapper.UserRoleToDB(domainUserRole)
	require.Equal(t, dbUserRole, backToDB)
}
