package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	"github.com/lechitz/AionApi/internal/record/core/usecase"
	tagdomain "github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Update_InvalidIDs(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	_, err := suite.RecordService.Update(suite.Ctx, 0, 1, input.UpdateRecordCommand{})
	require.ErrorIs(t, err, usecase.ErrInvalidRecordIDOrUserID)
}

func TestService_Update_GetByIDError(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	recordID := uint64(2)

	suite.RecordRepository.EXPECT().
		GetByID(gomock.Any(), recordID, userID).
		Return(domain.Record{}, errors.New("db error"))

	_, err := suite.RecordService.Update(suite.Ctx, recordID, userID, input.UpdateRecordCommand{})
	require.ErrorIs(t, err, usecase.ErrGetRecord)
}

func TestService_Update_TagNotFound(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	recordID := uint64(2)
	newTagID := uint64(10)

	existing := domain.Record{ID: recordID, UserID: userID, TagID: 1, EventTime: time.Now().UTC()}
	suite.RecordRepository.EXPECT().
		GetByID(gomock.Any(), recordID, userID).
		Return(existing, nil)

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), newTagID, userID).
		Return(tagdomain.Tag{}, errors.New("not found"))

	cmd := input.UpdateRecordCommand{TagID: &newTagID}
	_, err := suite.RecordService.Update(suite.Ctx, recordID, userID, cmd)
	require.ErrorIs(t, err, usecase.ErrUpdateRecord)
	assert.Contains(t, err.Error(), "tag not found")
}

func TestService_Update_RepositoryError(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	recordID := uint64(2)
	existing := domain.Record{ID: recordID, UserID: userID, TagID: 1, EventTime: time.Now().UTC()}

	suite.RecordRepository.EXPECT().
		GetByID(gomock.Any(), recordID, userID).
		Return(existing, nil)

	suite.RecordRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(domain.Record{}, errors.New("db error"))

	_, err := suite.RecordService.Update(suite.Ctx, recordID, userID, input.UpdateRecordCommand{})
	require.ErrorIs(t, err, usecase.ErrUpdateRecord)
}

func TestService_Update_Success(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	recordID := uint64(2)
	when := time.Date(2024, 1, 5, 15, 0, 0, 0, time.UTC)
	existing := domain.Record{ID: recordID, UserID: userID, TagID: 3, EventTime: when}
	updated := existing
	updated.Description = stringPtr("updated")

	suite.RecordRepository.EXPECT().
		GetByID(gomock.Any(), recordID, userID).
		Return(existing, nil)

	suite.RecordRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(updated, nil)

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), updated.TagID, userID).
		Return(tagdomain.Tag{ID: updated.TagID, CategoryID: 10}, nil)

	suite.RecordCache.EXPECT().DeleteRecord(gomock.Any(), updated.ID, updated.UserID).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByDay(gomock.Any(), updated.UserID, when.Truncate(24*time.Hour)).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByCategory(gomock.Any(), uint64(10), updated.UserID).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByTag(gomock.Any(), updated.TagID, updated.UserID).Return(nil)

	cmd := input.UpdateRecordCommand{Description: stringPtr("updated")}
	result, err := suite.RecordService.Update(suite.Ctx, recordID, userID, cmd)
	require.NoError(t, err)
	assert.Equal(t, updated, result)
}

func TestService_Delete_Errors(t *testing.T) {
	tests := []struct {
		name    string
		getErr  error
		delErr  error
		wantErr error
	}{
		{
			name:    "get by id error",
			getErr:  errors.New("db error"),
			wantErr: usecase.ErrGetRecord,
		},
		{
			name:    "delete error",
			delErr:  errors.New("delete error"),
			wantErr: usecase.ErrDeleteRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			userID := uint64(1)
			recordID := uint64(2)
			existing := domain.Record{ID: recordID, UserID: userID, TagID: 1, EventTime: time.Now().UTC()}

			suite.RecordRepository.EXPECT().
				GetByID(gomock.Any(), recordID, userID).
				Return(existing, tt.getErr)

			if tt.getErr == nil {
				suite.RecordRepository.EXPECT().
					Delete(gomock.Any(), recordID, userID).
					Return(tt.delErr)
			}

			err := suite.RecordService.Delete(suite.Ctx, recordID, userID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_Delete_Success(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	recordID := uint64(2)
	when := time.Date(2024, 1, 5, 15, 0, 0, 0, time.UTC)
	existing := domain.Record{ID: recordID, UserID: userID, TagID: 3, EventTime: when}

	suite.RecordRepository.EXPECT().
		GetByID(gomock.Any(), recordID, userID).
		Return(existing, nil)

	suite.RecordRepository.EXPECT().
		Delete(gomock.Any(), recordID, userID).
		Return(nil)

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), existing.TagID, userID).
		Return(tagdomain.Tag{ID: existing.TagID, CategoryID: 10}, nil)

	suite.RecordCache.EXPECT().DeleteRecord(gomock.Any(), recordID, userID).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByDay(gomock.Any(), userID, when.Truncate(24*time.Hour)).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByCategory(gomock.Any(), uint64(10), userID).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByTag(gomock.Any(), existing.TagID, userID).Return(nil)

	err := suite.RecordService.Delete(suite.Ctx, recordID, userID)
	require.NoError(t, err)
}

func TestService_DeleteAll(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	suite.RecordRepository.EXPECT().
		DeleteAllByUser(gomock.Any(), userID).
		Return(nil)

	err := suite.RecordService.DeleteAll(suite.Ctx, userID)
	require.NoError(t, err)
}

func TestService_ListRecordsByCategory(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	categoryID := uint64(2)
	records := []domain.Record{{ID: 1, UserID: userID}}

	suite.RecordRepository.EXPECT().
		ListByCategory(gomock.Any(), categoryID, userID, 10, nil, nil).
		Return(records, nil)

	result, err := suite.RecordService.ListRecordsByCategory(suite.Ctx, categoryID, userID, 10)
	require.NoError(t, err)
	assert.Equal(t, records, result)
}

func TestService_ListAllUntil_RepositoryError(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	until := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)

	suite.RecordRepository.EXPECT().
		ListAllUntil(gomock.Any(), userID, until, 10).
		Return(nil, errors.New("db error"))

	result, err := suite.RecordService.ListAllUntil(suite.Ctx, userID, until, 10)
	require.Error(t, err)
	assert.Contains(t, err.Error(), usecase.FailedToListRecords)
	assert.Nil(t, result)
}

func TestService_SearchRecords_Limits(t *testing.T) {
	tests := []struct {
		name          string
		inputLimit    int
		expectedLimit int
	}{
		{
			name:          "default limit",
			inputLimit:    0,
			expectedLimit: 20,
		},
		{
			name:          "cap limit",
			inputLimit:    150,
			expectedLimit: 100,
		},
		{
			name:          "use provided limit",
			inputLimit:    25,
			expectedLimit: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			userID := uint64(1)
			filters := domain.SearchFilters{Query: "hello", Limit: tt.inputLimit}
			records := []domain.Record{{ID: 1, UserID: userID}}

			suite.RecordRepository.EXPECT().
				SearchRecords(gomock.Any(), userID, gomock.AssignableToTypeOf(domain.SearchFilters{})).
				DoAndReturn(func(_ context.Context, _ uint64, f domain.SearchFilters) ([]domain.Record, error) {
					assert.Equal(t, tt.expectedLimit, f.Limit)
					return records, nil
				})

			result, err := suite.RecordService.SearchRecords(suite.Ctx, userID, filters)
			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}

func TestService_SearchRecords_Error(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	filters := domain.SearchFilters{Query: "oops", Limit: 10}

	suite.RecordRepository.EXPECT().
		SearchRecords(gomock.Any(), userID, filters).
		Return(nil, errors.New("search failed"))

	result, err := suite.RecordService.SearchRecords(suite.Ctx, userID, filters)
	require.Error(t, err)
	assert.Nil(t, result)
}
