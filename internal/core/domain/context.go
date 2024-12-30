package domain

import "context"

type ContextControl struct {
	Context         context.Context
	CancelCauseFunc context.CancelCauseFunc
}
