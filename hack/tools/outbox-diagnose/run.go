package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lechitz/AionApi/internal/adapter/secondary/db/postgres"
	outboxrepo "github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/repository"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	platformcfg "github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/keygen"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

func run(cfg config) error {
	if cfg.envFile != "" {
		if err := loadEnvFile(cfg.envFile); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	conf, err := loadAppConfig()
	if err != nil {
		return err
	}
	if conf.DB.Host == "postgres" {
		conf.DB.Host = "localhost"
	}

	dbConn, err := postgres.NewConnection(backgroundContext(), conf.DB, noopLogger{})
	if err != nil {
		return err
	}
	defer postgres.Close(dbConn, noopLogger{})

	repo := outboxrepo.NewEventRepository(postgres.NewDBAdapter(dbConn), noopLogger{})

	stats, err := repo.GetStats(backgroundContext())
	if err != nil {
		return err
	}
	pending, err := repo.ListByStatus(backgroundContext(), "pending", cfg.sampleLimit)
	if err != nil {
		return err
	}
	failed, err := repo.ListFailed(backgroundContext(), cfg.sampleLimit)
	if err != nil {
		return err
	}

	printStats(stats)
	printSample("pending sample", pending)
	printSample("failed sample", failed)
	return nil
}

func loadAppConfig() (*platformcfg.Config, error) {
	loader := platformcfg.New(staticKeyGen("hack-outbox-diagnose-secret"))
	return loader.Load(noopLogger{})
}

type staticKeyGen string

func (s staticKeyGen) Generate() (string, error) {
	return string(s), nil
}

type noopLogger struct{}

func (noopLogger) Infof(string, ...any)                      {}
func (noopLogger) Errorf(string, ...any)                     {}
func (noopLogger) Debugf(string, ...any)                     {}
func (noopLogger) Warnf(string, ...any)                      {}
func (noopLogger) Infow(string, ...any)                      {}
func (noopLogger) Errorw(string, ...any)                     {}
func (noopLogger) Debugw(string, ...any)                     {}
func (noopLogger) Warnw(string, ...any)                      {}
func (noopLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopLogger) DebugwCtx(context.Context, string, ...any) {}

func printStats(stats domain.Stats) {
	oldestPendingAge := "n/a"
	if stats.OldestPendingAtUTC != nil && !stats.OldestPendingAtUTC.IsZero() {
		oldestPendingAge = time.Since(stats.OldestPendingAtUTC.UTC()).Round(time.Second).String()
	}

	_, _ = fmt.Fprintf(os.Stdout, "pending_count=%d published_count=%d failed_count=%d oldest_pending_age=%s\n",
		stats.PendingCount,
		stats.PublishedCount,
		stats.FailedCount,
		oldestPendingAge,
	)
}

func printSample(title string, events []domain.Event) {
	_, _ = fmt.Fprintf(os.Stdout, "%s:\n", title)
	if len(events) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, "  (empty)")
		return
	}

	for _, event := range events {
		_, _ = fmt.Fprintf(os.Stdout, "  event_id=%s aggregate=%s/%s type=%s attempts=%d available_at=%s last_error=%q\n",
			event.EventID,
			event.AggregateType,
			event.AggregateID,
			event.EventType,
			event.AttemptCount,
			event.AvailableAtUTC.UTC().Format(time.RFC3339),
			event.LastError,
		)
	}
}

var (
	_ keygen.Generator     = staticKeyGen("")
	_ logger.ContextLogger = noopLogger{}
)
