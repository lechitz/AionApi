package kafka

import (
	"strings"

	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	kafkago "github.com/segmentio/kafka-go"
)

type ProjectionEventReader struct {
	reader *kafkago.Reader
	logger logger.ContextLogger
}

func NewProjectionEventReader(brokers string, groupID string, topic string, log logger.ContextLogger) *ProjectionEventReader {
	return &ProjectionEventReader{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:  splitBrokers(brokers),
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
		logger: log,
	}
}

func (r *ProjectionEventReader) Close() error {
	return r.reader.Close()
}

func splitBrokers(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}
