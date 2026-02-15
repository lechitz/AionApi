package main

import (
	"context"
	"errors"
	"os"
	"testing"
)

type fakeLifecycleApp struct {
	startErr    error
	stopErr     error
	doneCh      chan os.Signal
	startCalled bool
	stopCalled  bool
}

func (f *fakeLifecycleApp) Start(context.Context) error {
	f.startCalled = true
	return f.startErr
}

func (f *fakeLifecycleApp) Stop(context.Context) error {
	f.stopCalled = true
	return f.stopErr
}

func (f *fakeLifecycleApp) Done() <-chan os.Signal {
	return f.doneCh
}

func TestRunWithFactory(t *testing.T) {
	t.Run("start error returns 1", func(t *testing.T) {
		app := &fakeLifecycleApp{
			startErr: errors.New("start failed"),
			doneCh:   make(chan os.Signal),
		}

		code := runWithFactory(func() lifecycleApp { return app })
		if code != 1 {
			t.Fatalf("expected exit code 1, got %d", code)
		}
		if !app.startCalled {
			t.Fatal("expected start to be called")
		}
		if app.stopCalled {
			t.Fatal("did not expect stop to be called when start fails")
		}
	})

	t.Run("stop error returns 1", func(t *testing.T) {
		done := make(chan os.Signal)
		close(done)
		app := &fakeLifecycleApp{
			stopErr: errors.New("stop failed"),
			doneCh:  done,
		}

		code := runWithFactory(func() lifecycleApp { return app })
		if code != 1 {
			t.Fatalf("expected exit code 1, got %d", code)
		}
		if !app.startCalled || !app.stopCalled {
			t.Fatal("expected start and stop to be called")
		}
	})

	t.Run("success returns 0", func(t *testing.T) {
		done := make(chan os.Signal)
		close(done)
		app := &fakeLifecycleApp{doneCh: done}

		code := runWithFactory(func() lifecycleApp { return app })
		if code != 0 {
			t.Fatalf("expected exit code 0, got %d", code)
		}
		if !app.startCalled || !app.stopCalled {
			t.Fatal("expected start and stop to be called")
		}
	})
}
