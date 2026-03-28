package mapper_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserMapperRoundtrip(t *testing.T) {
	now := time.Now().UTC()
	deleted := now.Add(-time.Hour)
	locale := "pt-BR"

	dbUser := model.UserDB{
		ID:        1,
		Name:      "Le",
		Username:  "lechitz",
		Email:     "le@example.com",
		Password:  "hashed",
		Locale:    &locale,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: gorm.DeletedAt{Time: deleted, Valid: true},
	}

	domainUser := mapper.UserFromDB(dbUser)
	require.NotNil(t, domainUser.DeletedAt)
	require.Equal(t, deleted, *domainUser.DeletedAt)
	require.Equal(t, locale, *domainUser.Locale)

	back := mapper.UserToDB(domainUser)
	require.True(t, back.DeletedAt.Valid)
	require.Equal(t, deleted, back.DeletedAt.Time)

	back = mapper.UserToDB(domain.User{ID: 10, Username: "x"})
	require.False(t, back.DeletedAt.Valid)
}
