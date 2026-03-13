package main

import "flag"

const defaultSampleLimit = 5
const defaultEnvFile = "infrastructure/docker/environments/dev/.env.dev"

type config struct {
	sampleLimit int
	envFile     string
}

func loadConfig(args []string) config {
	fs := flag.NewFlagSet("outbox-diagnose", flag.ContinueOnError)
	cfg := config{}

	fs.IntVar(&cfg.sampleLimit, "sample-limit", defaultSampleLimit, "sample size for pending and failed outbox rows")
	fs.StringVar(&cfg.envFile, "env-file", defaultEnvFile, "optional env file loaded before reading app config")
	_ = fs.Parse(args)
	return cfg
}
