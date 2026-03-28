package main

import (
	"github.com/lechitz/aion-api/internal/platform/fxapp"
	"go.uber.org/fx"
)

func newFXApp() lifecycleApp {
	return fx.New(
		fxapp.InfraModule,
		fxapp.OutboxPublisherModule,
	)
}
