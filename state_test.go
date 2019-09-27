package go_fsm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ctxWithState(t *testing.T) {
	ctx := ctxWithState(context.Background(), "idle")
	assert.Equal(t, State("idle"), ctx.Value(stateCtxKey).(State))
}

func Test_StateFromCtx(t *testing.T) {
	t.Run("Valid context", func(t *testing.T) {
		ctx := ctxWithState(context.Background(), "idle")
		state := StateFromCtx(ctx)
		assert.Equal(t, "idle", state)
	})
}
