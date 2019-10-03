package go_fsm

import "context"

type (
	ctxEventKey int
	Event       = string
)

var eventCtxKey ctxEventKey

func ctxWithEvent(ctx context.Context, event Event) context.Context {
	return context.WithValue(ctx, eventCtxKey, event)
}

func EventFromCtx(ctx context.Context) (Event, error) {
	event, ok := ctx.Value(eventCtxKey).(Event)
	if !ok {
		return "", ErrCanNotExtractEvent
	}

	return event, nil
}
