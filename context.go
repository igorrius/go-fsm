package go_fsm

import "context"

// check context for nil and replace it with context.Background
func checkAndFixEmptyContext(ctx context.Context) context.Context {
	if nil == ctx {
		return context.Background()
	}
	return ctx
}
