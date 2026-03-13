package usecase

import (
	"context"
	"testing"

	eventoutboxdomain "github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	eventoutboxinput "github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	dbport "github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/record/core/ports/output"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

type recordRepositoryWithDBWrapper struct {
	*mocks.MockRecordRepository
	cloned output.RecordRepository
	gotDB  dbport.DB
}

func (s *recordRepositoryWithDBWrapper) WithDB(database dbport.DB) output.RecordRepository {
	s.gotDB = database
	return s.cloned
}

type outboxServiceWithDBStub struct {
	cloned eventoutboxinput.Service
	gotDB  dbport.DB
}

func (s *outboxServiceWithDBStub) Enqueue(context.Context, eventoutboxdomain.Event) error {
	return nil
}

func (s *outboxServiceWithDBStub) WithDB(database dbport.DB) eventoutboxinput.Service {
	s.gotDB = database
	return s.cloned
}

type outboxServiceNoop struct{}

func (s *outboxServiceNoop) Enqueue(context.Context, eventoutboxdomain.Event) error {
	return nil
}

func TestRunWithinRecordOutboxTransaction_UsesTransactionBoundDependencies(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txManager := mocks.NewMockDB(ctrl)
	txDB := mocks.NewMockDB(ctrl)
	baseRecordRepo := mocks.NewMockRecordRepository(ctrl)
	txRecordRepo := mocks.NewMockRecordRepository(ctrl)

	recordRepo := &recordRepositoryWithDBWrapper{
		MockRecordRepository: baseRecordRepo,
		cloned:               txRecordRepo,
	}
	outboxTxClone := &outboxServiceNoop{}
	outboxBase := &outboxServiceWithDBStub{cloned: outboxTxClone}

	txManager.EXPECT().WithContext(gomock.Any()).Return(txManager)
	txManager.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fn func(dbport.DB) error) error {
		return fn(txDB)
	})

	service := &Service{
		RecordRepository:   recordRepo,
		OutboxService:      outboxBase,
		TransactionManager: txManager,
	}

	var gotRecordRepo output.RecordRepository
	var gotOutbox eventoutboxinput.Service
	err := service.runWithinRecordOutboxTransaction(context.Background(), func(recordRepo output.RecordRepository, outboxService eventoutboxinput.Service) error {
		gotRecordRepo = recordRepo
		gotOutbox = outboxService
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if recordRepo.gotDB != txDB {
		t.Fatalf("expected record repo to receive transaction db")
	}
	if outboxBase.gotDB != txDB {
		t.Fatalf("expected outbox service to receive transaction db")
	}
	if gotRecordRepo != txRecordRepo {
		t.Fatalf("expected transaction-bound record repository")
	}
	if gotOutbox != outboxTxClone {
		t.Fatalf("expected transaction-bound outbox service")
	}
}

func TestRunWithinRecordOutboxTransaction_FallsBackWithoutTransactionManager(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseRecordRepo := mocks.NewMockRecordRepository(ctrl)
	baseOutbox := &outboxServiceNoop{}

	service := &Service{
		RecordRepository: baseRecordRepo,
		OutboxService:    baseOutbox,
	}

	var gotRecordRepo output.RecordRepository
	var gotOutbox eventoutboxinput.Service
	err := service.runWithinRecordOutboxTransaction(context.Background(), func(recordRepo output.RecordRepository, outboxService eventoutboxinput.Service) error {
		gotRecordRepo = recordRepo
		gotOutbox = outboxService
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotRecordRepo != baseRecordRepo {
		t.Fatalf("expected base record repository when transaction manager is nil")
	}
	if gotOutbox != baseOutbox {
		t.Fatalf("expected base outbox service when transaction manager is nil")
	}
}
