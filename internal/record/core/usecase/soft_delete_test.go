package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/usecase"
	tagdomain "github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
