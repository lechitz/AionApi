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

func TestCreate_InvalidTagID(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	in := gmodel.CreateRecordInput{TagID: "bad"}
	out, err := h.Create(t.Context(), in, 1)

	require.ErrorIs(t, err, controller.ErrInvalidRecordID)
	assert.Nil(t, out)
}

func TestCreate_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	in := gmodel.CreateRecordInput{TagID: "1"}
	out, err := h.Create(t.Context(), in, 0)

	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestCreate_ServiceError(t *testing.T) {
	expected := errors.New("create failed")
	svc := &recordServiceStub{
		createFn: func(_ context.Context, _ input.CreateRecordCommand) (domain.Record, error) {
			return domain.Record{}, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	in := gmodel.CreateRecordInput{TagID: "1"}
	out, err := h.Create(t.Context(), in, 1)

	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestCreate_Success(t *testing.T) {
	desc := "desc"
	eventTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	recordedAt := time.Date(2024, 1, 2, 4, 5, 6, 0, time.UTC)
	duration := int32(120)
	value := 12.5
	source := "mobile"
	tz := "UTC"
	status := "done"

	svc := &recordServiceStub{
		createFn: func(_ context.Context, cmd input.CreateRecordCommand) (domain.Record, error) {
			require.Equal(t, uint64(10), cmd.UserID)
			require.Equal(t, uint64(42), cmd.TagID)
			require.NotNil(t, cmd.Description)
			require.Equal(t, desc, *cmd.Description)
			require.Equal(t, eventTime, cmd.EventTime)
			require.NotNil(t, cmd.RecordedAt)
			require.Equal(t, recordedAt, *cmd.RecordedAt)
			require.NotNil(t, cmd.DurationSecs)
			require.Equal(t, int(duration), *cmd.DurationSecs)
			require.NotNil(t, cmd.Value)
			require.InEpsilon(t, value, *cmd.Value, 1e-6)
			require.NotNil(t, cmd.Source)
			require.Equal(t, source, *cmd.Source)
			require.NotNil(t, cmd.Timezone)
			require.Equal(t, tz, *cmd.Timezone)
			require.NotNil(t, cmd.Status)
			require.Equal(t, status, *cmd.Status)

			return domain.Record{
				ID:           99,
				UserID:       10,
				TagID:        42,
				Description:  &desc,
				EventTime:    eventTime,
				RecordedAt:   &recordedAt,
				DurationSecs: ptrInt(3600),
				Value:        &value,
				Source:       &source,
				Timezone:     &tz,
				Status:       &status,
				CreatedAt:    eventTime,
				UpdatedAt:    eventTime,
			}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	eventStr := eventTime.Format(time.RFC3339)
	recordedStr := recordedAt.Format(time.RFC3339)
	in := gmodel.CreateRecordInput{
		TagID:           "42",
		Description:     &desc,
		EventTime:       &eventStr,
		RecordedAt:      &recordedStr,
		DurationSeconds: &duration,
		Value:           &value,
		Source:          &source,
		Timezone:        &tz,
		Status:          &status,
	}

	out, err := h.Create(t.Context(), in, 10)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "99", out.ID)
	assert.Equal(t, "10", out.UserID)
	assert.Equal(t, "42", out.TagID)
	assert.Equal(t, desc, *out.Description)
	assert.Equal(t, eventStr, out.EventTime)
	require.NotNil(t, out.RecordedAt)
	assert.Equal(t, recordedStr, *out.RecordedAt)
	require.NotNil(t, out.DurationSeconds)
	assert.Equal(t, int32(3600), *out.DurationSeconds)
	assert.InEpsilon(t, value, *out.Value, 1e-6)
	assert.Equal(t, source, *out.Source)
	assert.Equal(t, tz, *out.Timezone)
	assert.Equal(t, status, *out.Status)
}

func ptrInt(v int) *int {
	return &v
}
