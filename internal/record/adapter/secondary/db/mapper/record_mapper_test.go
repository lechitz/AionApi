package mapper_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/stretchr/testify/require"
)

func TestRecordMapperRoundtripAndSlice(t *testing.T) {
	now := time.Now().UTC()
	recordedAt := now.Add(-time.Minute)
	duration := 15
	value := 3.14
	source := "manual"
	tz := "UTC"
	status := "active"
	desc := "desc"
	deletedAt := now.Add(-time.Hour)

	dbRecord := model.Record{
		ID:           1,
		UserID:       2,
		Description:  &desc,
		TagID:        3,
		EventTime:    now,
		RecordedAt:   &recordedAt,
		DurationSecs: &duration,
		Value:        &value,
		Source:       &source,
		Timezone:     &tz,
		Status:       &status,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    &deletedAt,
	}

	domainRecord := mapper.RecordFromDB(dbRecord)
	require.Equal(t, dbRecord.ID, domainRecord.ID)
	require.NotNil(t, domainRecord.Description)

	backToDB := mapper.RecordToDB(domainRecord)
	require.Equal(t, dbRecord.TagID, backToDB.TagID)

	list := mapper.RecordsFromDB([]model.Record{dbRecord})
	require.Len(t, list, 1)
	require.Equal(t, uint64(1), list[0].ID)
}
