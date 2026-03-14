package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"testing"
	"time"
)

type fakeLifecycleApp struct {
	startErr    error
	stopErr     error
	doneCh      chan os.Signal
	startCalled bool
	stopCalled  bool
	startCtx    context.Context
	stopCtx     context.Context
}

func (f *fakeLifecycleApp) Start(ctx context.Context) error {
	f.startCalled = true
	f.startCtx = ctx
	return f.startErr
}

func (f *fakeLifecycleApp) Stop(ctx context.Context) error {
	f.stopCalled = true
	f.stopCtx = ctx
	return f.stopErr
}

func (f *fakeLifecycleApp) Done() <-chan os.Signal {
	return f.doneCh
}

func testRunConfig() runConfig {
	return runConfig{
		startTimeout: 50 * time.Millisecond,
		stopTimeout:  50 * time.Millisecond,
		logf:         func(slog.Level, string, ...any) {},
	}
}

func assertWaitForShutdownReturns(t *testing.T, fn func()) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("waitForShutdownSignal did not return in time")
	}
}

func TestWaitForShutdownSignalReturnsOnAppDone(t *testing.T) {
	appDone := make(chan os.Signal)
	extraSignal := make(chan os.Signal)
	osSignal := make(chan struct{})
	close(appDone)

	assertWaitForShutdownReturns(t, func() {
		waitForShutdownSignal(appDone, extraSignal, osSignal)
	})
}

func TestWaitForShutdownSignalReturnsOnExtraSignal(t *testing.T) {
	appDone := make(chan os.Signal)
	extraSignal := make(chan os.Signal, 1)
	osSignal := make(chan struct{})
	extraSignal <- syscall.SIGTERM

	assertWaitForShutdownReturns(t, func() {
		waitForShutdownSignal(appDone, extraSignal, osSignal)
	})
}

func TestWaitForShutdownSignalReturnsOnOSSignal(t *testing.T) {
	appDone := make(chan os.Signal)
	extraSignal := make(chan os.Signal)
	osSignal := make(chan struct{})
	close(osSignal)

	assertWaitForShutdownReturns(t, func() {
		waitForShutdownSignal(appDone, extraSignal, osSignal)
	})
}

func TestRunWithFactoryStartErrorReturns1(t *testing.T) {
	app := &fakeLifecycleApp{
		startErr: errors.New("start failed"),
		doneCh:   make(chan os.Signal),
	}

	code := runWithFactory(func() lifecycleApp { return app }, testRunConfig())
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !app.startCalled {
		t.Fatal("expected start to be called")
	}
	if app.stopCalled {
		t.Fatal("did not expect stop to be called when start fails")
	}
}

func TestRunWithFactoryStopErrorReturns1(t *testing.T) {
	done := make(chan os.Signal)
	close(done)
	app := &fakeLifecycleApp{
		stopErr: errors.New("stop failed"),
		doneCh:  done,
	}

	code := runWithFactory(func() lifecycleApp { return app }, testRunConfig())
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithFactorySuccessReturns0(t *testing.T) {
	done := make(chan os.Signal)
	close(done)
	app := &fakeLifecycleApp{doneCh: done}

	code := runWithFactory(func() lifecycleApp { return app }, testRunConfig())
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithFactoryExternalSignalChannelTriggersShutdown(t *testing.T) {
	done := make(chan os.Signal)
	signalCh := make(chan os.Signal, 1)
	signalCh <- syscall.SIGTERM
	app := &fakeLifecycleApp{doneCh: done}
	cfg := testRunConfig()
	cfg.signalCh = signalCh

	code := runWithFactory(func() lifecycleApp { return app }, cfg)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithFactoryTimedSignalTriggersShutdown(t *testing.T) {
	done := make(chan os.Signal)
	signalCh := make(chan os.Signal, 1)
	app := &fakeLifecycleApp{doneCh: done}
	cfg := testRunConfig()
	cfg.signalCh = signalCh

	go func() {
		time.Sleep(10 * time.Millisecond)
		signalCh <- syscall.SIGTERM
	}()

	code := runWithFactory(func() lifecycleApp { return app }, cfg)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithFactoryReceivesDeadlineContexts(t *testing.T) {
	done := make(chan os.Signal)
	close(done)
	app := &fakeLifecycleApp{doneCh: done}
	cfg := testRunConfig()

	code := runWithFactory(func() lifecycleApp { return app }, cfg)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if app.startCtx == nil || app.stopCtx == nil {
		t.Fatal("expected start and stop contexts to be captured")
	}
	if _, ok := app.startCtx.Deadline(); !ok {
		t.Fatal("expected start context to have deadline")
	}
	if _, ok := app.stopCtx.Deadline(); !ok {
		t.Fatal("expected stop context to have deadline")
	}
}

func TestRunWithFactoryRealSignalTriggersShutdownInSubprocess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real signal subprocess test in short mode")
	}
	if runtime.GOOS == "windows" {
		t.Skip("SIGTERM behavior differs on windows")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()

	//nolint:gosec // Intentional helper-process pattern executing current test binary.
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestRunWithFactory_RealSignalProcess", "--")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	if err := cmd.Run(); err != nil {
		t.Fatalf("expected helper process to exit successfully, got %v", err)
	}
}

func TestRunWithFactory_RealSignalProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	done := make(chan os.Signal)
	app := &fakeLifecycleApp{doneCh: done}
	cfg := testRunConfig()
	cfg.signalCh = nil
	cfg.logf = func(slog.Level, string, ...any) {}

	go func() {
		time.Sleep(30 * time.Millisecond)
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	code := runWithFactory(func() lifecycleApp { return app }, cfg)
	if code != 0 || !app.startCalled || !app.stopCalled {
		os.Exit(1)
	}
	os.Exit(0)
}

func TestRunWithFactoryInvalidStartTimeoutReturns1(t *testing.T) {
	cfg := testRunConfig()
	cfg.startTimeout = 0
	factoryCalled := false

	code := runWithFactory(func() lifecycleApp {
		factoryCalled = true
		return &fakeLifecycleApp{doneCh: make(chan os.Signal)}
	}, cfg)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if factoryCalled {
		t.Fatal("did not expect factory to be called when start timeout is invalid")
	}
}

func TestRunWithFactoryInvalidStopTimeoutReturns1(t *testing.T) {
	cfg := testRunConfig()
	cfg.stopTimeout = 0
	factoryCalled := false

	code := runWithFactory(func() lifecycleApp {
		factoryCalled = true
		return &fakeLifecycleApp{doneCh: make(chan os.Signal)}
	}, cfg)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if factoryCalled {
		t.Fatal("did not expect factory to be called when stop timeout is invalid")
	}
}

func TestRunWithFactoryNilFactoryReturns1(t *testing.T) {
	code := runWithFactory(nil, testRunConfig())
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRunWithFactoryNilLogfFallbackStillRuns(t *testing.T) {
	done := make(chan os.Signal)
	close(done)
	app := &fakeLifecycleApp{doneCh: done}
	cfg := testRunConfig()
	cfg.logf = nil

	code := runWithFactory(func() lifecycleApp { return app }, cfg)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithDepsInvalidBootstrapEnvReturns1AndSkipsFactory(t *testing.T) {
	factoryCalled := false
	logCalled := false

	code := runWithDeps(
		func() lifecycleApp {
			factoryCalled = true
			return &fakeLifecycleApp{doneCh: make(chan os.Signal)}
		},
		func(key string) string {
			if key == envBootstrapStartTimeout {
				return "invalid"
			}
			return ""
		},
		func(slog.Level, string, ...any) {
			logCalled = true
		},
	)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if factoryCalled {
		t.Fatal("did not expect factory to be called on bootstrap config error")
	}
	if !logCalled {
		t.Fatal("expected error log to be called")
	}
}

func TestRunWithDepsValidBootstrapEnvRuns(t *testing.T) {
	done := make(chan os.Signal)
	close(done)
	app := &fakeLifecycleApp{doneCh: done}

	code := runWithDeps(
		func() lifecycleApp { return app },
		func(key string) string {
			switch key {
			case envBootstrapStartTimeout:
				return "2s"
			case envBootstrapStopTimeout:
				return "2s"
			default:
				return ""
			}
		},
		func(slog.Level, string, ...any) {},
	)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithDepsNilGetenvFallbackStillRuns(t *testing.T) {
	done := make(chan os.Signal)
	close(done)
	app := &fakeLifecycleApp{doneCh: done}

	code := runWithDeps(
		func() lifecycleApp { return app },
		nil,
		func(slog.Level, string, ...any) {},
	)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithDepsNilLogfFallbackStillRuns(t *testing.T) {
	done := make(chan os.Signal)
	close(done)
	app := &fakeLifecycleApp{doneCh: done}

	code := runWithDeps(
		func() lifecycleApp { return app },
		func(string) string { return "" },
		nil,
	)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !app.startCalled || !app.stopCalled {
		t.Fatal("expected start and stop to be called")
	}
}

func TestRunWithDepsNilFactoryReturns1(t *testing.T) {
	code := runWithDeps(
		nil,
		func(string) string { return "" },
		func(slog.Level, string, ...any) {},
	)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}
