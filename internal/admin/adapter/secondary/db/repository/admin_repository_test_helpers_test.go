package repository_test

import (
	"testing"

	repository "github.com/lechitz/aion-api/internal/admin/adapter/secondary/db/repository"
	"github.com/lechitz/aion-api/tests/mocks"
	"go.uber.org/mock/gomock"
)

func newAdminRepo(t *testing.T) (*repository.AdminRepository, *mocks.MockDB) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)
	loggerMock := mocks.NewMockContextLogger(ctrl)

	loggerMock.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	return repository.New(dbMock, loggerMock), dbMock
}
