package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type lifecycleApp interface {
	Start(context.Context) error
	Stop(context.Context) error
	Done() <-chan os.Signal
}

type runConfig struct {
	startTimeout time.Duration
	stopTimeout  time.Duration
	signalCh     <-chan os.Signal
	logf         func(slog.Level, string, ...any)
}

func defaultBootstrapLogf(level slog.Level, msg string, args ...any) {
	slog.Default().Log(context.Background(), level, msg, args...)
}

func runWithDeps(factory func() lifecycleApp, getenv func(string) string, logf func(slog.Level, string, ...any)) int {
	if logf == nil {
		logf = defaultBootstrapLogf
	}
	if getenv == nil {
		getenv = os.Getenv
	}

	cfg, err := loadBootstrapConfig(getenv)
	if err != nil {
		logf(slog.LevelError, "failed to load bootstrap config", "error", err)
		return 1
	}

	return runWithFactory(factory, runConfig{
		startTimeout: cfg.StartTimeout,
		stopTimeout:  cfg.StopTimeout,
		logf:         logf,
	})
}

func runWithFactory(factory func() lifecycleApp, cfg runConfig) int {
	if cfg.logf == nil {
		cfg.logf = defaultBootstrapLogf
	}
	if factory == nil || cfg.startTimeout <= 0 || cfg.stopTimeout <= 0 {
		cfg.logf(slog.LevelError, "invalid bootstrap runtime")
		return 1
	}

	app := factory()

	startCtx, cancelStart := context.WithTimeout(context.Background(), cfg.startTimeout)
	defer cancelStart()
	if err := app.Start(startCtx); err != nil {
		cfg.logf(slog.LevelError, "failed to start bootstrap app", "error", err)
		return 1
	}

	osSignalCtx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	waitForShutdownSignal(app.Done(), cfg.signalCh, osSignalCtx.Done())

	stopCtx, cancelStop := context.WithTimeout(context.Background(), cfg.stopTimeout)
	defer cancelStop()
	if err := app.Stop(stopCtx); err != nil {
		cfg.logf(slog.LevelError, "failed to stop bootstrap app", "error", err)
		return 1
	}

	return 0
}

func waitForShutdownSignal(appDone, extraSignal <-chan os.Signal, osSignal <-chan struct{}) {
	if extraSignal == nil {
		select {
		case <-appDone:
		case <-osSignal:
		}
		return
	}

	select {
	case <-appDone:
	case <-extraSignal:
	case <-osSignal:
	}
}
