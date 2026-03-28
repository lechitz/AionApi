package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	eventoutboxdomain "github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	tagdomain "github.com/lechitz/aion-api/internal/tag/core/domain"
	"github.com/lechitz/aion-api/tests/mocks"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create_Success(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	tagID := uint64(10)
	eventTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	description := "Morning workout"
	duration := 3600
	value := 150.0
	source := "mobile"

	ctx := context.WithValue(suite.Ctx, ctxkeys.UserID, userID)

	cmd := input.CreateRecordCommand{
		TagID:        tagID,
		Description:  &description,
		EventTime:    eventTime,
		DurationSecs: &duration,
		Value:        &value,
		Source:       &source,
	}

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), tagID, userID).
		Return(tagdomain.Tag{ID: tagID, Name: "Exercise", CategoryID: 1}, nil)

	expectedRecord := domain.Record{
		ID:           1,
		UserID:       userID,
		TagID:        tagID,
		Description:  &description,
		EventTime:    eventTime,
		DurationSecs: &duration,
		Value:        &value,
		Source:       &source,
		Status:       stringPtr(usecase.DefaultRecordStatus),
		Timezone:     stringPtr(usecase.DefaultTimezone),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	suite.RecordRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, rec domain.Record) (domain.Record, error) {
			rec.ID = expectedRecord.ID
			rec.CreatedAt = expectedRecord.CreatedAt
			rec.UpdatedAt = expectedRecord.UpdatedAt
			return rec, nil
		})

	suite.RecordCache.EXPECT().SaveRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByDay(gomock.Any(), userID, gomock.Any()).Return(nil).AnyTimes()
	suite.TagRepository.EXPECT().GetByID(gomock.Any(), tagID, userID).Return(tagdomain.Tag{ID: tagID, CategoryID: 1}, nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByCategory(gomock.Any(), gomock.Any(), userID).Return(nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByTag(gomock.Any(), tagID, userID).Return(nil).AnyTimes()

	result, err := suite.RecordService.Create(ctx, cmd)

	require.NoError(t, err)
	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, tagID, result.TagID)
	assert.Equal(t, description, *result.Description)
	assert.Equal(t, eventTime, result.EventTime)
	assert.NotNil(t, result.Status)
	assert.Equal(t, usecase.DefaultRecordStatus, *result.Status)
	assert.NotNil(t, result.Timezone)
	assert.Equal(t, usecase.DefaultTimezone, *result.Timezone)
}

func TestService_Create_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		setupCtx  func(t *testing.T) context.Context
		cmd       input.CreateRecordCommand
		setupMock func(*setup.RecordServiceTestSuite)
		wantErr   error
	}{
		{
			name:     "error - user not authenticated",
			setupCtx: func(t *testing.T) context.Context { return t.Context() },
			cmd: input.CreateRecordCommand{
				TagID:     10,
				EventTime: time.Now().UTC(),
			},
			setupMock: func(_ *setup.RecordServiceTestSuite) {},
			wantErr:   usecase.ErrUserNotAuthenticated,
		},
		{
			name: "error - recordedAt in future",
			setupCtx: func(t *testing.T) context.Context {
				return context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
			},
			cmd: input.CreateRecordCommand{
				TagID:      10,
				EventTime:  time.Now().UTC(),
				RecordedAt: timePtr(time.Now().UTC().Add(24 * time.Hour)),
			},
			setupMock: func(_ *setup.RecordServiceTestSuite) {},
			wantErr:   usecase.ErrRecordedAtFuture,
		},
		{
			name: "error - tag not found",
			setupCtx: func(t *testing.T) context.Context {
				return context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
			},
			cmd: input.CreateRecordCommand{
				TagID:     999,
				EventTime: time.Now().UTC(),
			},
			setupMock: func(s *setup.RecordServiceTestSuite) {
				s.TagRepository.EXPECT().
					GetByID(gomock.Any(), uint64(999), uint64(1)).
					Return(tagdomain.Tag{}, errors.New("tag not found"))
			},
			wantErr: usecase.ErrCreateRecord,
		},
		{
			name: "error - repository failure",
			setupCtx: func(t *testing.T) context.Context {
				return context.WithValue(t.Context(), ctxkeys.UserID, uint64(1))
			},
			cmd: input.CreateRecordCommand{
				TagID:     10,
				EventTime: time.Now().UTC(),
			},
			setupMock: func(s *setup.RecordServiceTestSuite) {
				s.TagRepository.EXPECT().
					GetByID(gomock.Any(), uint64(10), uint64(1)).
					Return(tagdomain.Tag{ID: 10}, nil)

				s.RecordRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(domain.Record{}, errors.New("database error"))
			},
			wantErr: usecase.ErrCreateRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			ctx := tt.setupCtx(t)
			tt.setupMock(suite)

			result, err := suite.RecordService.Create(ctx, tt.cmd)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr.Error())
			assert.Equal(t, domain.Record{}, result)
		})
	}
}

func TestService_Create_EnqueuesOutboxBestEffort(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	tagID := uint64(10)
	eventTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	ctx := context.WithValue(suite.Ctx, ctxkeys.UserID, userID)
	ctx = context.WithValue(ctx, ctxkeys.TraceID, "trace-123")
	ctx = context.WithValue(ctx, ctxkeys.RequestID, "req-123")

	outbox := mocks.NewMockOutboxService(suite.Ctrl)
	suite.RecordService.WithOutbox(outbox)

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), tagID, userID).
		Return(tagdomain.Tag{ID: tagID, Name: "Exercise", CategoryID: 1}, nil)

	suite.RecordRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, rec domain.Record) (domain.Record, error) {
			rec.ID = 77
			rec.CreatedAt = time.Now().UTC()
			rec.UpdatedAt = rec.CreatedAt
			return rec, nil
		})

	suite.RecordCache.EXPECT().SaveRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByDay(gomock.Any(), userID, gomock.Any()).Return(nil).AnyTimes()
	suite.TagRepository.EXPECT().GetByID(gomock.Any(), tagID, userID).Return(tagdomain.Tag{ID: tagID, CategoryID: 1}, nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByCategory(gomock.Any(), gomock.Any(), userID).Return(nil).AnyTimes()
	suite.RecordCache.EXPECT().DeleteRecordsByTag(gomock.Any(), tagID, userID).Return(nil).AnyTimes()

	outbox.EXPECT().Enqueue(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, event eventoutboxdomain.Event) error {
		assert.Equal(t, usecase.RecordAggregateType, event.AggregateType)
		assert.Equal(t, "77", event.AggregateID)
		assert.Equal(t, usecase.RecordEventTypeCreatedV1, event.EventType)
		assert.Equal(t, usecase.RecordEventVersionV1, event.EventVersion)
		assert.Equal(t, "trace-123", event.TraceID)
		assert.Equal(t, "req-123", event.RequestID)
		assert.NotEmpty(t, event.PayloadJSON)
		return errors.New("outbox unavailable")
	})

	result, err := suite.RecordService.Create(ctx, input.CreateRecordCommand{
		TagID:     tagID,
		EventTime: eventTime,
	})

	require.NoError(t, err)
	assert.Equal(t, uint64(77), result.ID)
}

// Helper functions.
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
