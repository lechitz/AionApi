package main

import (
	"flag"
	"time"
)

const (
	defaultRealtimeSmokeHost         = "http://localhost:5001"
	defaultRealtimeSmokeUser         = "testuser"
	defaultRealtimeSmokePassword     = "Test@123"
	defaultRealtimeSmokeTagID        = "24"
	defaultRealtimeSmokeTimeout      = 20 * time.Second
	defaultRealtimeSmokePollInterval = 500 * time.Millisecond
)

type config struct {
	host         string
	username     string
	password     string
	tagID        string
	timeout      time.Duration
	pollInterval time.Duration
}

func loadConfig(args []string) config {
	fs := flag.NewFlagSet("realtime-record-smoke", flag.ContinueOnError)
	cfg := config{}

	fs.StringVar(&cfg.host, "host", defaultRealtimeSmokeHost, "AionApi host base URL")
	fs.StringVar(&cfg.username, "username", defaultRealtimeSmokeUser, "login username")
	fs.StringVar(&cfg.password, "password", defaultRealtimeSmokePassword, "login password")
	fs.StringVar(&cfg.tagID, "tag-id", defaultRealtimeSmokeTagID, "record tag ID used in smoke flow")
	fs.DurationVar(&cfg.timeout, "timeout", defaultRealtimeSmokeTimeout, "overall smoke timeout")
	fs.DurationVar(&cfg.pollInterval, "poll-interval", defaultRealtimeSmokePollInterval, "projection polling interval")

	_ = fs.Parse(args)
	return cfg
}
