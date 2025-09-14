// Package runtime provides a lightweight runtime group to start and stop multiple actors together.
package runtime

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// Log message constants.
const (
	// LogShutdownSignal is logged when the context is canceled and shutdown begins.
	LogShutdownSignal = "shutdown signal received"
	// LogUnexpectedFailure is logged when an actor returns an unexpected error.
	LogUnexpectedFailure = "unexpected server failure"
)

// Actor represents a unit of work with start and interrupt functions.
type Actor struct {
	Start     func() error
	Interrupt func(error)
}

// Group manages a set of actors that are started and stopped together.
type Group struct {
	actors []Actor
}

// Add registers an actor to the group.
func (g *Group) Add(start func() error, interrupt func(error)) {
	g.actors = append(g.actors, Actor{Start: start, Interrupt: interrupt})
}

// Run starts all actors and blocks until the context is canceled or any actor fails.
// When exiting, all actors receive an interrupt signal.
func (g *Group) Run(ctx context.Context, _ time.Duration, log logger.ContextLogger) {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, a := range g.actors {
		wg.Add(1)
		go func(a Actor) {
			defer wg.Done()
			if err := a.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				select {
				case errCh <- err:
				default:
				}
			}
		}(a)
	}

	select {
	case <-ctx.Done():
		log.Infow(LogShutdownSignal)
	case err := <-errCh:
		log.Errorw(LogUnexpectedFailure, commonkeys.Error, err.Error())
	}

	var iwg sync.WaitGroup
	for _, a := range g.actors {
		iwg.Add(1)
		go func(a Actor) {
			defer iwg.Done()
			a.Interrupt(context.Canceled)
		}(a)
	}
	iwg.Wait()
	wg.Wait()
}
