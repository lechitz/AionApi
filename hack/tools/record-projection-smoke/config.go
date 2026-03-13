package main

import (
	"flag"
	"time"
)

const (
	defaultSmokeHost     = "http://localhost:5001"
	defaultSmokeUser     = "testuser"
	defaultSmokePassword = "Test@123"
	defaultSmokeTagID    = "24"
	defaultSmokeTimeout  = 20 * time.Second
	defaultSmokePoll     = 500 * time.Millisecond
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
	fs := flag.NewFlagSet("record-projection-smoke", flag.ContinueOnError)
	cfg := config{}

	fs.StringVar(&cfg.host, "host", defaultSmokeHost, "AionApi host base URL")
	fs.StringVar(&cfg.username, "username", defaultSmokeUser, "login username")
	fs.StringVar(&cfg.password, "password", defaultSmokePassword, "login password")
	fs.StringVar(&cfg.tagID, "tag-id", defaultSmokeTagID, "record tag ID used in smoke flow")
	fs.DurationVar(&cfg.timeout, "timeout", defaultSmokeTimeout, "overall smoke timeout")
	fs.DurationVar(&cfg.pollInterval, "poll-interval", defaultSmokePoll, "projection polling interval")

	_ = fs.Parse(args)
	return cfg
}
