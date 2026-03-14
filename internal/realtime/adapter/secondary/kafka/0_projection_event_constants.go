// Package kafka reads projection-ready events from Kafka for realtime fan-out.
package kafka

const (
	// LogRealtimeConsumerStarted is emitted when the realtime Kafka loop starts.
	LogRealtimeConsumerStarted = "realtime projection consumer started"
	// LogRealtimeConsumerStopped is emitted when the realtime Kafka loop stops.
	LogRealtimeConsumerStopped = "realtime projection consumer stopped"
	// LogRealtimeConsumeFailed is emitted when reading from Kafka fails.
	LogRealtimeConsumeFailed = "realtime projection consume failed"
)
