package mapper_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/admin/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/admin/core/domain"
	userdomain "github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/stretchr/testify/require"
)

func TestAdminUserFromUserAndBack(t *testing.T) {
	now := time.Now().UTC()
	deleted := now.Add(-time.Hour)

	usr := userdomain.User{
		ID:        10,
		Username:  "jane",
		Email:     "jane@example.com",
		Name:      "Jane",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: &deleted,
	}

	adminUser := mapper.AdminUserFromUser(usr)
	require.Equal(t, usr.ID, adminUser.ID)
	require.False(t, adminUser.IsActive)

	backToUser := mapper.UserFromAdminUser(domain.AdminUser{
		ID:        11,
		Username:  "john",
		Email:     "john@example.com",
		Name:      "John",
		CreatedAt: now,
		UpdatedAt: now,
	})
	require.Equal(t, uint64(11), backToUser.ID)
	require.Equal(t, "john", backToUser.Username)
	require.Nil(t, backToUser.DeletedAt)
}
