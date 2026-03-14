package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/adapter/primary/graphql/controller"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.Update(t.Context(), gmodel.UpdateRecordInput{ID: "1"}, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestUpdate_InvalidRecordID(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.Update(t.Context(), gmodel.UpdateRecordInput{ID: "bad"}, 1)
	require.ErrorIs(t, err, controller.ErrInvalidRecordID)
	assert.Nil(t, out)
}

func TestUpdate_ServiceError(t *testing.T) {
	expected := errors.New("update failed")
	svc := &recordServiceStub{
		updateFn: func(_ context.Context, _, _ uint64, _ input.UpdateRecordCommand) (domain.Record, error) {
			return domain.Record{}, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.Update(t.Context(), gmodel.UpdateRecordInput{ID: "10"}, 1)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestUpdate_Success(t *testing.T) {
	desc := "updated"
	value := 9.5
	source := "web"
	tz := "UTC"
	status := "ok"
	tagID := "7"
	eventTime := time.Date(2024, 1, 5, 12, 0, 0, 0, time.UTC)
	recordedAt := time.Date(2024, 1, 5, 12, 30, 0, 0, time.UTC)
	duration := int32(33)

	svc := &recordServiceStub{
		updateFn: func(_ context.Context, recordID uint64, userID uint64, cmd input.UpdateRecordCommand) (domain.Record, error) {
			require.Equal(t, uint64(10), recordID)
			require.Equal(t, uint64(5), userID)
			require.NotNil(t, cmd.TagID)
			require.Equal(t, uint64(7), *cmd.TagID)
			require.NotNil(t, cmd.EventTime)
			require.Equal(t, eventTime, *cmd.EventTime)
			require.NotNil(t, cmd.RecordedAt)
			require.Equal(t, recordedAt, *cmd.RecordedAt)
			require.NotNil(t, cmd.DurationSecs)
			require.Equal(t, int(duration), *cmd.DurationSecs)
			require.Equal(t, desc, *cmd.Description)
			require.InEpsilon(t, value, *cmd.Value, 1e-6)
			require.Equal(t, source, *cmd.Source)
			require.Equal(t, tz, *cmd.Timezone)
			require.Equal(t, status, *cmd.Status)

			return domain.Record{
				ID:        recordID,
				UserID:    userID,
				TagID:     7,
				EventTime: eventTime,
				CreatedAt: eventTime,
				UpdatedAt: eventTime,
			}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	eventStr := eventTime.Format(time.RFC3339)
	recordedStr := recordedAt.Format(time.RFC3339)
	in := gmodel.UpdateRecordInput{
		ID:              "10",
		Description:     &desc,
		TagID:           &tagID,
		EventTime:       &eventStr,
		RecordedAt:      &recordedStr,
		DurationSeconds: &duration,
		Value:           &value,
		Source:          &source,
		Timezone:        &tz,
		Status:          &status,
	}

	out, err := h.Update(t.Context(), in, 5)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "10", out.ID)
	assert.Equal(t, "5", out.UserID)
	assert.Equal(t, "7", out.TagID)
}

func TestUpdate_IgnoresInvalidOptionalFields(t *testing.T) {
	badTag := "bad"
	badTime := "bad"

	svc := &recordServiceStub{
		updateFn: func(_ context.Context, _ uint64, _ uint64, cmd input.UpdateRecordCommand) (domain.Record, error) {
			require.Nil(t, cmd.TagID)
			require.Nil(t, cmd.EventTime)
			require.Nil(t, cmd.RecordedAt)
			return domain.Record{ID: 1, UserID: 2, TagID: 3, EventTime: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	in := gmodel.UpdateRecordInput{ID: "1", TagID: &badTag, EventTime: &badTime, RecordedAt: &badTime}
	_, err := h.Update(t.Context(), in, 2)
	require.NoError(t, err)
}
