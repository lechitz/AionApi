package usecase

import (
	"context"
	"testing"

	eventoutboxdomain "github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
	eventoutboxinput "github.com/lechitz/aion-api/internal/eventoutbox/core/ports/input"
	dbport "github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/record/core/ports/output"
	"github.com/lechitz/aion-api/tests/mocks"
	"go.uber.org/mock/gomock"
)

type txAwareRecordRepository struct {
	output.RecordRepository
	txRepository output.RecordRepository
}

func (r txAwareRecordRepository) WithDB(_ dbport.DB) output.RecordRepository {
	return r.txRepository
}

type txAwareOutboxService struct {
	eventoutboxinput.Service
	txService eventoutboxinput.Service
}

func (s txAwareOutboxService) WithDB(_ dbport.DB) eventoutboxinput.Service {
	return s.txService
}

type noopOutboxServiceWithTx struct{}

func (noopOutboxServiceWithTx) Enqueue(_ context.Context, _ eventoutboxdomain.Event) error {
	return nil
}

func TestService_RunWithinRecordOutboxTransaction_UsesTransactionBoundDependencies(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbMock := mocks.NewMockDB(ctrl)
	txMock := mocks.NewMockDB(ctrl)
	recordRepo := mocks.NewMockRecordRepository(ctrl)
	txRecordRepo := mocks.NewMockRecordRepository(ctrl)

	service := &Service{
		RecordRepository:   txAwareRecordRepository{RecordRepository: recordRepo, txRepository: txRecordRepo},
		OutboxService:      txAwareOutboxService{Service: noopOutboxServiceWithTx{}, txService: noopOutboxServiceWithTx{}},
		TransactionManager: dbMock,
	}

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fn func(dbport.DB) error) error {
		return fn(txMock)
	})

	var gotRepo output.RecordRepository
	var gotOutbox eventoutboxinput.Service
	err := service.runWithinRecordOutboxTransaction(t.Context(), func(recordRepo output.RecordRepository, outboxService eventoutboxinput.Service) error {
		gotRepo = recordRepo
		gotOutbox = outboxService
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotRepo != txRecordRepo {
		t.Fatalf("expected transaction-bound record repo")
	}
	if gotOutbox == nil {
		t.Fatalf("expected transaction-bound outbox service")
	}
}
