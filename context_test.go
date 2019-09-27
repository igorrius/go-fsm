package go_fsm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_checkAndFixEmptyContext(t *testing.T) {
	t.Run("Valid context", func(t *testing.T) {
		ctxA := context.WithValue(context.Background(), "testKey", "testValue")
		ctxB := checkAndFixEmptyContext(ctxA)
		assert.Equal(t, ctxA, ctxB)
	})

	t.Run("Invalid context", func(t *testing.T) {
		ctx := checkAndFixEmptyContext(nil)
		assert.Equal(t, context.Background(), ctx)
	})
}
