package usecase_test

import (
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
	suite.RecordCache.EXPECT().DeleteRecordsByDay(gomock.Any(), updated.UserID, time.Date(when.UTC().Year(), when.UTC().Month(), when.UTC().Day(), 0, 0, 0, 0, time.UTC)).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByCategory(gomock.Any(), uint64(10), updated.UserID).Return(nil)
	suite.RecordCache.EXPECT().DeleteRecordsByTag(gomock.Any(), updated.TagID, updated.UserID).Return(nil)

	cmd := input.UpdateRecordCommand{Description: stringPtr("updated")}
	result, err := suite.RecordService.Update(suite.Ctx, recordID, userID, cmd)
	require.NoError(t, err)
	assert.Equal(t, updated, result)
}
