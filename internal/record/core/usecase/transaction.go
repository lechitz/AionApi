package usecase

import (
	"context"
	"errors"

	eventoutboxinput "github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	dbport "github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/record/core/ports/output"
)

type recordRepositoryWithDB interface {
	WithDB(database dbport.DB) output.RecordRepository
}

type outboxServiceWithDB interface {
	WithDB(database dbport.DB) eventoutboxinput.Service
}

func (s *Service) runWithinRecordOutboxTransaction(ctx context.Context, fn func(recordRepo output.RecordRepository, outboxService eventoutboxinput.Service) error) error {
	if s.TransactionManager == nil || s.OutboxService == nil {
		return fn(s.RecordRepository, s.OutboxService)
	}

	recordRepoFactory, ok := s.RecordRepository.(recordRepositoryWithDB)
	if !ok {
		return fn(s.RecordRepository, s.OutboxService)
	}

	outboxFactory, ok := s.OutboxService.(outboxServiceWithDB)
	if !ok {
		return fn(s.RecordRepository, s.OutboxService)
	}

	return s.TransactionManager.WithContext(ctx).Transaction(func(tx dbport.DB) error {
		return fn(recordRepoFactory.WithDB(tx), outboxFactory.WithDB(tx))
	})
}

func isGetRecordError(err error) bool {
	return errors.Is(err, ErrGetRecord)
}
