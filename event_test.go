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
	t.Run("Context with event value", func(t *testing.T) {
		ctx := ctxWithEvent(context.Background(), "event1")
		event, err := EventFromCtx(ctx)
		assert.Nil(t, err)
		assert.Equal(t, "event1", event)
	})

	t.Run("Context without event value", func(t *testing.T) {
		event, err := EventFromCtx(context.TODO())
		assert.EqualError(t, err, ErrCanNotExtractEvent.Error())
		assert.Equal(t, "", event)
	})
}
