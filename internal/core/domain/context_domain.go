package domain

import (
	"context"
	"time"
)

type CtxControl struct {
	Ctx       context.Context
	UserID    uint64
	Token     string
	RequestID string
	StartTime time.Time
}
