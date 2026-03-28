// Package fxapp wires the application using Uber Fx modules.
package fxapp

import (
	"context"
	"sync"
	"time"

	eventOutboxRepo "github.com/lechitz/aion-api/internal/eventoutbox/adapter/secondary/db/repository"
	eventOutboxKafka "github.com/lechitz/aion-api/internal/eventoutbox/adapter/secondary/kafka"
	eventOutboxInput "github.com/lechitz/aion-api/internal/eventoutbox/core/ports/input"
	eventOutboxOutput "github.com/lechitz/aion-api/internal/eventoutbox/core/ports/output"
	eventOutbox "github.com/lechitz/aion-api/internal/eventoutbox/core/usecase"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.uber.org/fx"
)

// OutboxPublisherModule wires the dedicated background publisher process.
//
//nolint:gochecknoglobals // Fx modules are intended as package-level options.
var OutboxPublisherModule = fx.Options(
	fx.Provide(
		ProvideOutboxEventRepository,
		ProvideKafkaEventPublisher,
		ProvideOutboxPublisherService,
	),
	fx.Invoke(RunOutboxPublisher),
)

type outboxPublisherParams struct {
	fx.In

	Cfg       *config.Config
	Log       logger.ContextLogger
	Repo      eventOutboxOutput.EventRepository
	Publisher eventOutboxOutput.EventPublisher
}

// ProvideOutboxEventRepository exposes the durable outbox repository for the publisher process.
func ProvideOutboxEventRepository(database db.DB, log logger.ContextLogger) eventOutboxOutput.EventRepository {
	return eventOutboxRepo.NewEventRepository(database, log)
}

// ProvideKafkaEventPublisher exposes the Kafka-backed outbox publisher and registers its cleanup.
func ProvideKafkaEventPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	log logger.ContextLogger,
) eventOutboxOutput.EventPublisher {
	publisher := eventOutboxKafka.NewEventPublisher(cfg.Kafka, log)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return publisher.Close()
		},
	})
	return publisher
}

// ProvideOutboxPublisherService creates the batch publication use case.
func ProvideOutboxPublisherService(params outboxPublisherParams) eventOutboxInput.PublisherService {
	return eventOutbox.NewPublisherService(params.Repo, params.Publisher, params.Log)
}

// RunOutboxPublisher starts the periodic background loop for Kafka publication.
func RunOutboxPublisher(
	lc fx.Lifecycle,
	cfg *config.Config,
	service eventOutboxInput.PublisherService,
	log logger.ContextLogger,
) {
	if !cfg.Outbox.PublishEnabled {
		log.Infow("outbox publisher disabled by configuration")
		return
	}

	var (
		wg     sync.WaitGroup
		cancel context.CancelFunc
	)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// #nosec G118 -- Cancel is stored here and invoked during Fx OnStop.
			workerCtx, workerCancel := context.WithCancel(context.Background())
			cancel = workerCancel
			wg.Add(1)

			go func() {
				defer wg.Done()
				ticker := time.NewTicker(cfg.Outbox.PublishInterval)
				defer ticker.Stop()

				for {
					if err := service.PublishPending(workerCtx, cfg.Outbox.BatchSize); err != nil && workerCtx.Err() == nil {
						log.ErrorwCtx(workerCtx, "outbox publish cycle failed",
							commonkeys.Error, err.Error(),
							"batch_size", cfg.Outbox.BatchSize,
						)
					}

					select {
					case <-workerCtx.Done():
						return
					case <-ticker.C:
					}
				}
			}()

			log.Infow("outbox publisher started",
				"publish_interval", cfg.Outbox.PublishInterval.String(),
				"batch_size", cfg.Outbox.BatchSize,
			)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}

			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case <-done:
				log.Infow("outbox publisher stopped")
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})
}
