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

func TestListByDay_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByDay(t.Context(), 0, "2024-01-01")
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestListByDay_InvalidDate(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByDay(t.Context(), 1, "bad")
	require.Error(t, err)
	assert.Equal(t, "invalid date format, expected YYYY-MM-DD or RFC3339", err.Error())
	assert.Nil(t, out)
}

func TestListByDay_ServiceError(t *testing.T) {
	expected := errors.New("list failed")
	svc := &recordServiceStub{
		listByDayFn: func(_ context.Context, _ uint64, _ time.Time) ([]domain.Record, error) {
			return nil, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByDay(t.Context(), 1, "2024-01-01")
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestListByDay_Success(t *testing.T) {
	when := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	svc := &recordServiceStub{
		listByDayFn: func(_ context.Context, userID uint64, date time.Time) ([]domain.Record, error) {
			require.Equal(t, uint64(9), userID)
			require.Equal(t, when, date)
			return []domain.Record{{ID: 1, UserID: userID, TagID: 2, EventTime: when, CreatedAt: when, UpdatedAt: when}}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByDay(t.Context(), 9, "2024-01-01")
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "1", out[0].ID)
}
