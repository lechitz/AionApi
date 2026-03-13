package main

import (
	"flag"
	"time"
)

const (
	defaultIngestHost      = "http://localhost:8091"
	defaultKafkaBroker     = "localhost:19092"
	defaultIngestTopic     = "aion.ingest.events.v1"
	defaultIngestSource    = "wearable.garmin"
	defaultIngestTimeout   = 20 * time.Second
	defaultIngestPollSleep = 250 * time.Millisecond
)

type config struct {
	ingestHost  string
	kafkaBroker string
	topic       string
	source      string
	timeout     time.Duration
	pollSleep   time.Duration
}

func loadConfig(args []string) config {
	fs := flag.NewFlagSet("ingest-event-smoke", flag.ContinueOnError)
	cfg := config{}

	fs.StringVar(&cfg.ingestHost, "ingest-host", defaultIngestHost, "aion-ingest base URL")
	fs.StringVar(&cfg.kafkaBroker, "kafka-broker", defaultKafkaBroker, "Kafka/Redpanda broker address")
	fs.StringVar(&cfg.topic, "topic", defaultIngestTopic, "ingest topic name")
	fs.StringVar(&cfg.source, "source", defaultIngestSource, "X-Aion-Source value")
	fs.DurationVar(&cfg.timeout, "timeout", defaultIngestTimeout, "overall smoke timeout")
	fs.DurationVar(&cfg.pollSleep, "poll-sleep", defaultIngestPollSleep, "retry sleep between polling attempts")

	_ = fs.Parse(args)
	return cfg
}
