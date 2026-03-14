package fxapp

import (
	"context"
	"os"
	"sync"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	realtimeKafka "github.com/lechitz/AionApi/internal/realtime/adapter/secondary/kafka"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.uber.org/fx"
)

// RealtimeModule wires the realtime event consumer into the application lifecycle.
//
//nolint:gochecknoglobals // Fx modules are declared as package-level options across the application wiring.
var RealtimeModule = fx.Options(
	fx.Provide(ProvideRealtimeProjectionReader),
	fx.Invoke(RunRealtimeProjectionConsumer),
)

// ProvideRealtimeProjectionReader builds the Kafka reader used by the realtime consumer.
func ProvideRealtimeProjectionReader(
	lc fx.Lifecycle,
	cfg *config.Config,
	log logger.ContextLogger,
) *realtimeKafka.ProjectionEventReader {
	groupID := cfg.Realtime.ConsumerGroupPrefix
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		groupID += "-" + hostname
	}

	reader := realtimeKafka.NewProjectionEventReader(
		cfg.Kafka.Brokers,
		groupID,
		cfg.Kafka.RecordProjectionEventsTopic,
		log,
	)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return reader.Close()
		},
	})
	return reader
}

// RunRealtimeProjectionConsumer starts the realtime consumer loop when the feature is enabled.
func RunRealtimeProjectionConsumer(
	lc fx.Lifecycle,
	cfg *config.Config,
	reader *realtimeKafka.ProjectionEventReader,
	deps *AppDependencies,
	log logger.ContextLogger,
) {
	if !cfg.Realtime.Enabled {
		log.Infow("realtime stream disabled by configuration")
		return
	}
	if deps == nil || deps.RealtimeService == nil {
		log.Warnw("realtime stream enabled but realtime service is unavailable")
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
				log.Infow(realtimeKafka.LogRealtimeConsumerStarted,
					"kafka_topic", cfg.Kafka.RecordProjectionEventsTopic,
				)

				for {
					event, err := reader.Read(workerCtx)
					if err != nil {
						if workerCtx.Err() != nil {
							return
						}
						log.ErrorwCtx(workerCtx, realtimeKafka.LogRealtimeConsumeFailed, commonkeys.Error, err.Error())
						continue
					}
					deps.RealtimeService.Publish(workerCtx, event)
				}
			}()
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
				log.Infow(realtimeKafka.LogRealtimeConsumerStopped)
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})
}
