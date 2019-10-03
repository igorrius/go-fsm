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

func StateFromCtx(ctx context.Context) (State, error) {
	state, ok := ctx.Value(stateCtxKey).(State)
	if !ok {
		return "", ErrCanNotExtractState
	}

	return state, nil
}
