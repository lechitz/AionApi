package usecase_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	eventoutboxdomain "github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	tagdomain "github.com/lechitz/aion-api/internal/tag/core/domain"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type captureOutboxService struct {
	events []eventoutboxdomain.Event
}

func (c *captureOutboxService) Enqueue(_ context.Context, event eventoutboxdomain.Event) error {
	c.events = append(c.events, event)
	return nil
}

func TestService_Create_EnqueuesOutboxEvent(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	outbox := &captureOutboxService{}
	suite.RecordService.WithOutbox(outbox)

	userID := uint64(7)
	tagID := uint64(10)
	eventTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	description := "Morning workout"
	source := "mobile"

	ctx := context.WithValue(suite.Ctx, ctxkeys.UserID, userID)
	ctx = context.WithValue(ctx, ctxkeys.TraceID, "trace-1")
	ctx = context.WithValue(ctx, ctxkeys.RequestID, "req-1")

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), tagID, userID).
		Return(tagdomain.Tag{ID: tagID, Name: "Exercise", CategoryID: 1}, nil)

	suite.RecordRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, rec domain.Record) (domain.Record, error) {
			rec.ID = 55
			rec.CreatedAt = time.Now().UTC()
			rec.UpdatedAt = rec.CreatedAt
			return rec, nil
		})

	suite.RecordCache.EXPECT().SaveRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByDay(gomock.Any(), userID, gomock.Any()).Return(nil).AnyTimes()
	suite.TagRepository.EXPECT().GetByID(gomock.Any(), tagID, userID).Return(tagdomain.Tag{ID: tagID, CategoryID: 1}, nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByCategory(gomock.Any(), gomock.Any(), userID).Return(nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByTag(gomock.Any(), tagID, userID).Return(nil).AnyTimes()

	_, err := suite.RecordService.Create(ctx, input.CreateRecordCommand{
		TagID:       tagID,
		Description: &description,
		EventTime:   eventTime,
		Source:      &source,
	})
	require.NoError(t, err)
	require.Len(t, outbox.events, 1)
	require.Equal(t, usecase.RecordAggregateType, outbox.events[0].AggregateType)
	require.Equal(t, "55", outbox.events[0].AggregateID)
	require.Equal(t, usecase.RecordEventTypeCreatedV1, outbox.events[0].EventType)
	require.Equal(t, usecase.RecordEventVersionV1, outbox.events[0].EventVersion)
	require.Equal(t, "trace-1", outbox.events[0].TraceID)
	require.Equal(t, "req-1", outbox.events[0].RequestID)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(outbox.events[0].PayloadJSON, &payload))
	require.EqualValues(t, 55, payload["record_id"])
	require.EqualValues(t, userID, payload["user_id"])
	require.EqualValues(t, tagID, payload["tag_id"])
	require.Equal(t, source, payload["source"])
}
