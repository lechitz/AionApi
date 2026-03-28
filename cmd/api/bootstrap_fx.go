package main

import (
	"github.com/lechitz/aion-api/internal/platform/fxapp"
	"go.uber.org/fx"
)

func newFXApp() lifecycleApp {
	return newFXAppWithOptions()
}

func newFXAppWithOptions(extraOptions ...fx.Option) lifecycleApp {
	options := []fx.Option{
		fx.Invoke(configureSwagger),
		fxapp.InfraModule,
		fxapp.ApplicationModule,
		fxapp.RealtimeModule,
		fxapp.ServerModule,
	}
	options = append(options, extraOptions...)

	return fx.New(fx.Options(options...))
}
