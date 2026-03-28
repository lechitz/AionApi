package repository_test

import (
	"testing"
	"time"

	repository "github.com/lechitz/aion-api/internal/record/adapter/secondary/db/repository"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/tests/mocks"
	"go.uber.org/mock/gomock"
)

func newRecordRepo(t *testing.T) (*repository.RecordRepository, *mocks.MockDB) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	dbMock := mocks.NewMockDB(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	return repository.New(dbMock, logger), dbMock
}

func newRecordRepoWithLogger(t *testing.T) (*repository.RecordRepository, *mocks.MockDB, *mocks.MockContextLogger) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	dbMock := mocks.NewMockDB(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	return repository.New(dbMock, logger), dbMock, logger
}

func sampleRecord() domain.Record {
	now := time.Now().UTC()
	desc := "desc"
	return domain.Record{ID: 1, UserID: 10, TagID: 20, Description: &desc, EventTime: now, CreatedAt: now, UpdatedAt: now}
}
