package mapper_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/tag/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTagMapperRoundtripAndSlice(t *testing.T) {
	now := time.Now().UTC()
	deleted := now.Add(-time.Hour)
	lastUsed := now.Add(-time.Minute)

	dbTag := model.TagDB{
		ID:          9,
		UserID:      7,
		CategoryID:  3,
		Name:        "Focus",
		Description: "high focus",
		Icon:        "bolt",
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   gorm.DeletedAt{Time: deleted, Valid: true},
		UsageCount:  5,
		LastUsedAt:  &lastUsed,
	}

	domainTag := mapper.TagFromDB(dbTag)
	require.NotNil(t, domainTag.DeletedAt)
	require.Equal(t, deleted, *domainTag.DeletedAt)

	mappedList := mapper.TagsFromDB([]model.TagDB{dbTag})
	require.Len(t, mappedList, 1)
	require.Equal(t, dbTag.Name, mappedList[0].Name)

	backToDB := mapper.TagToDB(mappedList[0])
	require.Equal(t, dbTag.ID, backToDB.ID)
	require.False(t, backToDB.DeletedAt.Valid)
}
