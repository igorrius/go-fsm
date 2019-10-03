package go_fsm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"sync/atomic"
	"testing"
)

func emptyStateActionFunc(nextState State) ActionFunc {
	return func(eventCtx EventContext, fsmCtx FsmContext) (next State, nextFsmCtx FsmContext, err error) {
		return nextState, nil, nil
	}
}

type customLogger struct {
}

func (l customLogger) Log(v ...interface{}) {
	log.Print(v...)
}

func (l customLogger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func TestNewFsm(t *testing.T) {
	t.Run("With default logger", func(t *testing.T) {
		fsm := NewFsm()
		assert.NotNil(t, fsm)
		assert.Equal(t, &nilLoggerAdapter{}, fsm.logger)
	})

	t.Run("With the custom logger", func(t *testing.T) {
		logger := &customLogger{}
		fsm := NewFsm(LoggerOption(logger))
		assert.NotNil(t, fsm)
		assert.Equal(t, logger, fsm.logger)
	})
}

func TestFsm_Init(t *testing.T) {
	t.Run("Init without existing state", func(t *testing.T) {
		_, err := NewFsm().InitWithState("idle")
		assert.EqualError(t, err, "invalid initial state [idle]")
	})

	t.Run("Init with existing state", func(t *testing.T) {
		fsm, err := NewFsm().
			When("someState", emptyStateActionFunc("idle")).
			InitWithState("someState")
		assert.NoError(t, err)
		assert.Equal(t, "someState", fsm.state)
	})
}

func TestFsm_When(t *testing.T) {
	t.Run("Empty action", func(t *testing.T) {
		fsm, err := NewFsm().When("idle", nil).InitWithState("idle")
		assert.NoError(t, err)
		if e := fsm.ProcessEvent("someEvent", nil); e != nil {
			assert.EqualError(t, e, ErrActionNotFound.Error())
			assert.Equal(t, "idle", fsm.state)
		}
	})

	t.Run("Invalid next state in action handler", func(t *testing.T) {
		fsm, err :=
			NewFsm().
				When("idle", emptyStateActionFunc("idle")).
				InitWithState("idle")
		assert.NoError(t, err)
		if e := fsm.ProcessEvent("someEvent", nil); e != nil {
			assert.EqualError(t, e, ErrActionNotFound.Error())
			assert.Equal(t, "idle", fsm.state)
		}
	})

	t.Run("Event from ctx", func(t *testing.T) {
		fsm, err :=
			NewFsm().
				When(
					"idle",
					func(eventCtx EventContext, fsmCtx FsmContext) (next State, nextFsmCtx FsmContext, err error) {
						event, err := EventFromCtx(eventCtx)
						assert.Nil(t, err)
						assert.Equal(t, "someEvent", event)
						return event, nil, nil
					}).
				InitWithState("idle")
		assert.NoError(t, err)
		if e := fsm.ProcessEvent("someEvent", nil); e != nil {
			assert.EqualError(t, e, ErrActionNotFound.Error())
			assert.Equal(t, "idle", fsm.state)
		}
	})

	t.Run("Valid init, state, event, action chain", func(t *testing.T) {
		fsm, err :=
			NewFsm().
				When("idle", emptyStateActionFunc("idle2")).
				When("idle2", emptyStateActionFunc("nextState")).
				When("nextState", emptyStateActionFunc("idle")).
				InitWithState("idle")
		assert.NoError(t, err)

		err = fsm.ProcessEvent("someEvent", context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, "idle2", fsm.state)

		err = fsm.ProcessEvent("someEvent", context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, "nextState", fsm.state)

		err = fsm.ProcessEvent("someEvent", context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, "idle", fsm.state)
	})

}

func TestFsm_Close(t *testing.T) {
	fsm, err :=
		NewFsm().
			When("idle", emptyStateActionFunc("idle")).
			InitWithState("idle")
	assert.NoError(t, err)

	err = fsm.ProcessEvent("someEvent", nil)
	assert.NoError(t, err)

	fsm.Close()
	err = fsm.ProcessEvent("someEvent", nil)
	assert.EqualError(t, err, "context canceled")
}

func TestFsm_Reset(t *testing.T) {
	fsm, err :=
		NewFsm().
			When("someInitState", emptyStateActionFunc("next")).
			When("next", emptyStateActionFunc("someInitState")).
			InitWithState("someInitState")
	assert.NoError(t, err)

	err = fsm.ProcessEvent("someEvent", nil)
	assert.NoError(t, err)
	assert.Equal(t, "next", fsm.state)

	assert.NoError(t, fsm.Reset())
	assert.Equal(t, "someInitState", fsm.state)
}

func TestFsm_AddTransitionFunc(t *testing.T) {
	t.Run("Post transition functions", func(t *testing.T) {
		var atomicCounter int32
		atomic.StoreInt32(&atomicCounter, 0)

		fsm, err :=
			NewFsm().
				When("idle", emptyStateActionFunc("nextState")).
				When("nextState", emptyStateActionFunc("idle")).
				InitWithState("idle")
		assert.NoError(t, err)

		fsm.RegisterPostTransitionFunc(
			"idle", "nextState",
			func(from, to State, fsmCtx FsmContext) error {
				atomic.AddInt32(&atomicCounter, 1)
				state, err := StateFromCtx(fsmCtx)
				assert.Nil(t, err)
				assert.Equal(t, state, from)
				assert.Equal(t, "idle", from)
				assert.Equal(t, "nextState", to)
				return nil
			})

		fsm.RegisterPostTransitionFunc(
			"idle", "nextState",
			func(from, to State, fsmCtx FsmContext) error {
				atomic.AddInt32(&atomicCounter, 1)
				state, err := StateFromCtx(fsmCtx)
				assert.Nil(t, err)
				assert.Equal(t, state, from)
				assert.Equal(t, "idle", from)
				assert.Equal(t, "nextState", to)
				return nil
			})

		fsm.RegisterPostTransitionFunc(
			"idle", "*",
			func(from, to State, fsmCtx FsmContext) error {
				atomic.AddInt32(&atomicCounter, 1)
				state, err := StateFromCtx(fsmCtx)
				assert.Nil(t, err)
				assert.Equal(t, state, from)
				assert.Equal(t, "idle", from)
				assert.Equal(t, "nextState", to)
				return nil
			})

		fsm.RegisterPostTransitionFunc(
			"*", "nextState",
			func(from, to State, fsmCtx FsmContext) error {
				atomic.AddInt32(&atomicCounter, 1)
				state, err := StateFromCtx(fsmCtx)
				assert.Nil(t, err)
				assert.Equal(t, state, from)
				assert.Equal(t, "idle", from)
				assert.Equal(t, "nextState", to)
				return nil
			})

		fsm.RegisterPostTransitionFunc(
			"*", "*",
			func(from, to State, fsmCtx FsmContext) error {
				atomic.AddInt32(&atomicCounter, 1)
				state, err := StateFromCtx(fsmCtx)
				assert.Nil(t, err)
				assert.Equal(t, state, from)
				assert.Equal(t, "idle", from)
				assert.Equal(t, "nextState", to)
				return nil
			})

		fsm.RegisterPostTransitionFunc(
			"nextState", "idle",
			func(from, to State, fsmCtx FsmContext) error {
				assert.Fail(t, "This function should not call")
				return nil
			})

		err = fsm.ProcessEvent("someEvent", nil)
		assert.NoError(t, err)
		assert.Equal(t, int32(5), atomic.LoadInt32(&atomicCounter))
	})
}

func TestFsm_isStateExists(t *testing.T) {
	fsm := &Fsm{
		actionMap: map[State]ActionFunc{
			"state1": func(eventCtx EventContext, fsmCtx FsmContext) (next State, nextFsmCtx FsmContext, err error) {
				return "state2", nil, nil
			},
			"state2": func(eventCtx EventContext, fsmCtx FsmContext) (next State, nextFsmCtx FsmContext, err error) {
				return "state1", nil, nil
			},
		},
	}

	assert.True(t, fsm.isStateExists("state1"))
	assert.True(t, fsm.isStateExists("state2"))
	assert.False(t, fsm.isStateExists("state3"))
}

func TestFsm_CurrentState(t *testing.T) {
	fsm := &Fsm{
		state: "idle",
	}
	assert.Equal(t, "idle", fsm.CurrentState())
}

func BenchmarkFsm_When(b *testing.B) {
	b.Run("Transition Permitted", func(b *testing.B) {
		b.ReportAllocs()

		fsm, err := NewFsm().
			When("idle", emptyStateActionFunc("next")).
			When("next", emptyStateActionFunc("idle")).
			InitWithState("idle")

		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fsm.ProcessEvent("benchmark", context.Background())
		}
	})

	b.Run("Transition Denied", func(b *testing.B) {
		b.ReportAllocs()

		fsm, err := NewFsm().
			When("idle", emptyStateActionFunc("next")).
			InitWithState("idle")

		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fsm.ProcessEvent("benchmark", context.Background())
		}
	})

	b.Run("Transition with post transition functions run", func(b *testing.B) {
		b.ReportAllocs()

		fsm, err := NewFsm().
			When("idle", emptyStateActionFunc("next")).
			When("next", emptyStateActionFunc("idle")).
			RegisterPostTransitionFunc(
				"idle", "next",
				func(from, to State, fsmCtx FsmContext) error {
					return nil
				}).
			RegisterPostTransitionFunc(
				"idle", "next",
				func(from, to State, fsmCtx FsmContext) error {
					return nil
				}).
			RegisterPostTransitionFunc(
				"next", "idle",
				func(from, to State, fsmCtx FsmContext) error {
					return nil
				}).
			RegisterPostTransitionFunc(
				"next", "idle",
				func(from, to State, fsmCtx FsmContext) error {
					return nil
				}).
			InitWithState("idle")

		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fsm.ProcessEvent("benchmark", context.Background())
		}
	})
}

func BenchmarkFsm_Reset(b *testing.B) {
	b.ReportAllocs()
	fsm, err := NewFsm().
		When("idle", emptyStateActionFunc("idle")).
		InitWithState("idle")

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fsm.Reset()
	}
}

func BenchmarkFsm_Close(b *testing.B) {
	b.Run("Close fsm", func(b *testing.B) {
		b.ReportAllocs()
		fsm, err := NewFsm().
			When("idle", emptyStateActionFunc("idle")).
			InitWithState("idle")

		if err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fsm.Close()
		}
	})

	b.Run("Close fsm and try to process event", func(b *testing.B) {
		b.ReportAllocs()
		fsm, err := NewFsm().
			When("idle", emptyStateActionFunc("idle")).
			InitWithState("idle")

		if err != nil {
			b.Fatal(err)
		}

		fsm.Close()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fsm.ProcessEvent("benchmark", context.Background())
		}
	})
}
