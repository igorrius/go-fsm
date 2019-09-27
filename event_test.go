package go_fsm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ctxWithEvent(t *testing.T) {
	ctx := ctxWithEvent(context.Background(), "event1")
	assert.Equal(t, Event("event1"), ctx.Value(eventCtxKey).(State))
}

func Test_EventFromCtx(t *testing.T) {
	t.Run("Valid context", func(t *testing.T) {
		ctx := ctxWithState(context.Background(), "event1")
		event := StateFromCtx(ctx)
		assert.Equal(t, "event1", event)
	})
}
