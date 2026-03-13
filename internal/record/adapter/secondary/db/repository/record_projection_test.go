package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRecordProjectionQueries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbMock := mocks.NewMockDB(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	repo := New(dbMock, logger)
	now := time.Date(2026, time.March, 13, 15, 0, 0, 0, time.UTC)

	t.Run("get projected by id success", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			row, ok := dest.(*recordProjectionRow)
			require.True(t, ok)
			row.RecordID = 5177
			row.UserID = 7
			row.TagID = 32
			row.LastEventID = "evt-1"
			row.LastEventType = "record.created"
			row.LastEventVersion = "v1"
			row.LastKafkaTopic = "aion.record.events.v1"
			row.LastKafkaPartition = 0
			row.LastKafkaOffset = 1
			row.EventTimeUTC = now
			row.LastConsumedAtUTC = now
			row.PayloadJSON = []byte(`{"record_id":5177}`)
			row.CreatedAtUTC = now
			row.UpdatedAtUTC = now
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.GetProjectedByID(t.Context(), 7, 5177)
		require.NoError(t, err)
		require.Equal(t, uint64(5177), got.RecordID)
		require.Equal(t, int64(1), got.LastKafkaOffset)
	})

	t.Run("list projected latest error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("query fail"))

		_, err := repo.ListProjectedLatest(t.Context(), 7, 5)
		require.Error(t, err)
	})

	t.Run("list projected page with cursor", func(t *testing.T) {
		at := now.Format(time.RFC3339)
		afterID := int64(5176)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Order(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Limit(10).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			rows, ok := dest.(*[]recordProjectionRow)
			require.True(t, ok)
			*rows = []recordProjectionRow{{
				RecordID:          5177,
				UserID:            7,
				TagID:             32,
				LastEventID:       "evt-2",
				LastEventType:     "record.updated",
				LastEventVersion:  "v1",
				LastKafkaTopic:    "aion.record.events.v1",
				LastKafkaOffset:   2,
				EventTimeUTC:      now,
				LastConsumedAtUTC: now,
				PayloadJSON:       []byte(`{"record_id":5177}`),
				CreatedAtUTC:      now,
				UpdatedAtUTC:      now,
			}}
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.ListProjectedPage(t.Context(), 7, 10, &at, &afterID)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, uint64(5177), got[0].RecordID)
	})
}
