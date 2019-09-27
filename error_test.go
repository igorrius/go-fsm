package go_fsm

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_checkErrors(t *testing.T) {
	t.Run("First error in arguments", func(t *testing.T) {
		errs := []error{errors.New("first error"), errors.New("second error")}
		err := checkErrors(errs...)
		assert.EqualError(t, err, "first error")
	})

	t.Run("Second error in arguments", func(t *testing.T) {
		errs := []error{nil, errors.New("second error")}
		err := checkErrors(errs...)
		assert.EqualError(t, err, "second error")
	})
}
