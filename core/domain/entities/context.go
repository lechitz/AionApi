package entities

import "context"

type ContextControl struct {
	BaseContext     context.Context
	CancelCauseFunc context.CancelCauseFunc
}
