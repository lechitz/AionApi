package main

import (
	"testing"

	swagger "github.com/lechitz/AionApi/contracts/openapi"
	"github.com/lechitz/AionApi/internal/platform/config"
)

func TestConfigureSwagger(t *testing.T) {
	cfg := &config.Config{
		General: config.GeneralConfig{
			Version: "v-test",
		},
		ServerHTTP: config.ServerHTTP{
			Context: "/aion/",
			APIRoot: "/api/v1/",
		},
	}

	configureSwagger(cfg)

	if swagger.SwaggerInfo.Title != SwaggerTitle {
		t.Fatalf("unexpected title: %s", swagger.SwaggerInfo.Title)
	}
	if swagger.SwaggerInfo.Version != "v-test" {
		t.Fatalf("unexpected version: %s", swagger.SwaggerInfo.Version)
	}
	if swagger.SwaggerInfo.BasePath != "/aion/api/v1" {
		t.Fatalf("unexpected base path: %s", swagger.SwaggerInfo.BasePath)
	}
}
