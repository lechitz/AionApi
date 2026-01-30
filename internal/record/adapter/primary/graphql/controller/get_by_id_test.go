package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/primary/graphql/controller"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByID_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.GetByID(t.Context(), 1, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestGetByID_RecordIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.GetByID(t.Context(), 0, 1)
	require.ErrorIs(t, err, controller.ErrRecordNotFound)
	assert.Nil(t, out)
}

func TestGetByID_ServiceError(t *testing.T) {
	expected := errors.New("not found")
	svc := &recordServiceStub{
		getByIDFn: func(_ context.Context, _, _ uint64) (domain.Record, error) {
			return domain.Record{}, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.GetByID(t.Context(), 1, 2)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestGetByID_Success_MapsDurationBounds(t *testing.T) {
	desc := "desc"
	overflow := 2147483648
	eventTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	svc := &recordServiceStub{
		getByIDFn: func(_ context.Context, _, _ uint64) (domain.Record, error) {
			return domain.Record{
				ID:           5,
				UserID:       6,
				TagID:        7,
				Description:  &desc,
				EventTime:    eventTime,
				DurationSecs: &overflow,
				CreatedAt:    eventTime,
				UpdatedAt:    eventTime,
			}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.GetByID(t.Context(), 5, 6)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "5", out.ID)
	assert.Equal(t, "6", out.UserID)
	assert.Equal(t, "7", out.TagID)
	assert.Equal(t, desc, *out.Description)
	assert.Nil(t, out.DurationSeconds)
}
