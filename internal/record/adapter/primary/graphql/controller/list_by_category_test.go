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

func TestListByCategory_UserIDMissing(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByCategory(t.Context(), 1, 0, 10)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)
	assert.Nil(t, out)
}

func TestListByCategory_CategoryIDInvalid(t *testing.T) {
	svc := &recordServiceStub{}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByCategory(t.Context(), 0, 1, 10)
	require.Error(t, err)
	assert.Equal(t, "category id cannot be zero", err.Error())
	assert.Nil(t, out)
}

func TestListByCategory_ServiceError(t *testing.T) {
	expected := errors.New("list failed")
	svc := &recordServiceStub{
		listByCatFn: func(_ context.Context, _, _ uint64, _ int) ([]domain.Record, error) {
			return nil, expected
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByCategory(t.Context(), 2, 3, 10)
	require.ErrorIs(t, err, expected)
	assert.Nil(t, out)
}

func TestListByCategory_Success(t *testing.T) {
	when := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	svc := &recordServiceStub{
		listByCatFn: func(_ context.Context, categoryID uint64, userID uint64, limit int) ([]domain.Record, error) {
			require.Equal(t, uint64(5), categoryID)
			require.Equal(t, uint64(6), userID)
			require.Equal(t, 3, limit)
			return []domain.Record{{ID: 1, UserID: userID, TagID: 9, EventTime: when, CreatedAt: when, UpdatedAt: when}}, nil
		},
	}
	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	out, err := h.ListByCategory(t.Context(), 5, 6, 3)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "1", out[0].ID)
}
