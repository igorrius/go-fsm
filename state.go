package go_fsm

import "context"

type (
	ctxStateKey int
	State       = string
)

var stateCtxKey ctxStateKey

func ctxWithState(ctx context.Context, state State) context.Context {
	return context.WithValue(ctx, stateCtxKey, state)
}

func StateFromCtx(ctx context.Context) State {
	state, ok := ctx.Value(stateCtxKey).(State)
	if !ok {
		panic("unknown state in context")
	}

	return state
}
