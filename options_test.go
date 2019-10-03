package go_fsm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newOptions(t *testing.T) {
	adapter := &nilLoggerAdapter{}
	options := newOptions(LoggerOption(adapter))
	assert.Equal(t, adapter, options.Logger)
}
