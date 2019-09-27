package go_fsm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newTransitionKey(t *testing.T) {
	key := newTransitionKey("from", "to")
	assert.Equal(t, transitionKey{"from", "to"}, key)
}
