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
	t.Run("Context with state", func(t *testing.T) {
		ctx := ctxWithState(context.Background(), "idle")
		state, err := StateFromCtx(ctx)
		assert.Nil(t, err)
		assert.Equal(t, "idle", state)
	})

	t.Run("Context without state", func(t *testing.T) {
		state, err := StateFromCtx(context.TODO())
		assert.EqualError(t, err, ErrCanNotExtractState.Error())
		assert.Equal(t, "", state)
	})
}
