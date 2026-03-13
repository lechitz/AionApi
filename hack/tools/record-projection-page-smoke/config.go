package main

import (
	"flag"
	"time"
)

const (
	defaultPageSmokeHost     = "http://localhost:5001"
	defaultPageSmokeUser     = "testuser"
	defaultPageSmokePassword = "Test@123"
	defaultPageSmokeTagID    = "24"
	defaultPageSmokeTimeout  = 25 * time.Second
	defaultPageSmokePoll     = 500 * time.Millisecond
	defaultPageSmokeLimit    = 2
)

type config struct {
	host         string
	username     string
	password     string
	tagID        string
	timeout      time.Duration
	pollInterval time.Duration
	pageLimit    int
}

func loadConfig(args []string) config {
	fs := flag.NewFlagSet("record-projection-page-smoke", flag.ContinueOnError)
	cfg := config{}

	fs.StringVar(&cfg.host, "host", defaultPageSmokeHost, "AionApi host base URL")
	fs.StringVar(&cfg.username, "username", defaultPageSmokeUser, "login username")
	fs.StringVar(&cfg.password, "password", defaultPageSmokePassword, "login password")
	fs.StringVar(&cfg.tagID, "tag-id", defaultPageSmokeTagID, "record tag ID used in smoke flow")
	fs.DurationVar(&cfg.timeout, "timeout", defaultPageSmokeTimeout, "overall smoke timeout")
	fs.DurationVar(&cfg.pollInterval, "poll-interval", defaultPageSmokePoll, "projection polling interval")
	fs.IntVar(&cfg.pageLimit, "page-limit", defaultPageSmokeLimit, "projection page size used in the smoke")

	_ = fs.Parse(args)
	return cfg
}
