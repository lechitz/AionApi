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

type Actor struct {
	Start     func() error
	Interrupt func(error)
}

type Group struct {
	actors []Actor
}

func (g *Group) Add(start func() error, interrupt func(error)) {
	g.actors = append(g.actors, Actor{Start: start, Interrupt: interrupt})
}

func (g *Group) Run(ctx context.Context, timeout time.Duration, log logger.ContextLogger) {
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
		log.Infow("shutdown signal received")
	case err := <-errCh:
		log.Errorw("unexpected server failure", commonkeys.Error, err.Error())
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
