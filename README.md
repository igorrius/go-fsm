[![Go Report Card](https://goreportcard.com/badge/github.com/igorrius/go-fsm)](https://goreportcard.com/report/github.com/igorrius/go-fsm)

# go-fsm
Finite State Machine in Go (with Blackjack and... using context)

Overview
----------
Finite State Machine is designed with alignment to
[Erlang Finite State Machine](http://erlang.org/documentation/doc-4.8.2/doc/design_principles/fsm.html) principles 
and inspired by https://github.com/dyrkin/fsm solution. But has a bit different approach which based on `context.Context`
GO package. 

A FSM can be described as a set of relations of the form:
```
State(S) x Event(E) -> Actions (A), State(S')
```
If we are in state S and the event E occurs, we should perform the actions A and make a transition to the state S'.

Install
-------
```
go get -u github.com/igorrius/go-fsm
```

Example
--------
```go
package main

import (
	"context"
	go_fsm "github.com/igorrius/go-fsm"
	"log"
)

type Logger struct {
}

func (l Logger) Log(v ...interface{}) {
	log.Print(v...)
}

func (l Logger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func main() {
	const (
		// define states
		stateIdle     = "idle"
		stateInAction = "inAction"
	)

	// create logger
	logger := &Logger{}

	// create a new instance of FSM,
	// if a logger is not needed use constructor without a logger option e.g. fsm, err := go_fsm.NewFsm().
	fsm, err := go_fsm.NewFsm(go_fsm.LoggerOption(logger)).
		When(
			stateIdle,
			func(eventCtx go_fsm.EventContext, fsmCtx go_fsm.FsmContext) (next go_fsm.State, nextFsmCtx go_fsm.FsmContext, err error) {
				var state go_fsm.State
				var event go_fsm.Event
				// get current state from eventCtx
				if event, err = go_fsm.EventFromCtx(eventCtx); err != nil {
					return
				}

				// get current state from fsmContext
				if state, err = go_fsm.StateFromCtx(fsmCtx); err != nil {
					return
				}

				switch event {
				case "moveRight":
					log.Println("Action: move right")
					nextFsmCtx = context.WithValue(fsmCtx, "degrees", "30")
					next = stateInAction
				case "moveLeft":
					log.Println("Action: move left")
					nextFsmCtx = context.WithValue(fsmCtx, "degrees", "50")
					next = stateInAction
				case "letsError":
					log.Println("Action: letsError")
					next = "unknownState"
				default:
					// FSM must stay in current state
					log.Println("Unknown event: ", event)
					next = state
				}

				return
			},
		).
		When(
			stateInAction,
			func(eventCtx go_fsm.EventContext, fsmCtx go_fsm.FsmContext) (next go_fsm.State, nextFsmCtx go_fsm.FsmContext, err error) {
				var state go_fsm.State
				var event go_fsm.Event
				// get current state from eventCtx
				if event, err = go_fsm.EventFromCtx(eventCtx); err != nil {
					return
				}

				// get current state from fsmContext
				if state, err = go_fsm.StateFromCtx(fsmCtx); err != nil {
					return
				}

				switch event {
				case "stop":
					log.Println("Action: stop")
					log.Println("Degrees: ", fsmCtx.Value("degrees").(string))
					next = stateIdle
				default:
					// FSM must stay in current state
					log.Println("Unknown event: ", event)
					next = state
				}

				return
			},
		).
		RegisterPostTransitionFunc("*", "*",
			func(from, to go_fsm.State, fsmCtx go_fsm.FsmContext) error {
				log.Println("Transition Function [ANY to ANY]")
				return nil
			},
		).
		RegisterPostTransitionFunc(stateIdle, "*",
			func(from, to go_fsm.State, fsmCtx go_fsm.FsmContext) error {
				log.Println("Transition Function [stateIdle to ANY]")
				return nil
			},
		).
		RegisterPostTransitionFunc(stateInAction, "*",
			func(from, to go_fsm.State, fsmCtx go_fsm.FsmContext) error {
				log.Println("Transition Function [stateInAction to ANY]")
				return nil
			},
		).
		RegisterPostTransitionFunc("*", stateInAction,
			func(from, to go_fsm.State, fsmCtx go_fsm.FsmContext) error {
				log.Println("Transition Function [ANY to stateInAction]")
				return nil
			},
		).
		RegisterPostTransitionFunc("*", stateIdle,
			func(from, to go_fsm.State, fsmCtx go_fsm.FsmContext) error {
				log.Println("Transition Function [ANY to stateIdle]")
				return nil
			},
		).
		RegisterPostTransitionFunc(stateInAction, stateIdle,
			func(from, to go_fsm.State, fsmCtx go_fsm.FsmContext) error {
				log.Println("Transition Function [stateIdle to stateInAction]")
				return nil
			},
		).
		RegisterPostTransitionFunc(stateIdle, stateInAction,
			func(from, to go_fsm.State, fsmCtx go_fsm.FsmContext) error {
				log.Println("Transition Function [stateIdle to stateInAction]")
				return nil
			},
		).
		InitWithState(stateIdle)

	if err != nil {
		log.Fatalln("FSM init error:", err)
	}
	defer fsm.Close()

	// call stop - wrong event
	if err = fsm.ProcessEvent("stop", context.TODO()); err != nil {
		log.Fatalln(err)
	}

	// call letsError - error next state
	if err = fsm.ProcessEvent("letsError", context.TODO()); err != nil {
		log.Println("This will be an error:", err)
	}

	// call move right event
	if err = fsm.ProcessEvent("moveRight", context.TODO()); err != nil {
		log.Fatalln(err)
	}

	// call move stop
	if err = fsm.ProcessEvent("stop", context.TODO()); err != nil {
		log.Fatalln(err)
	}

	// call move right event
	if err = fsm.ProcessEvent("moveLeft", context.TODO()); err != nil {
		log.Fatalln(err)
	}

	// call move stop - wrong event
	if err = fsm.ProcessEvent("stop", context.TODO()); err != nil {
		log.Println("This will be an error:", err)
	}
}

```

Benchmark
---------
```
Transition_Permitted-8                                   	   10000	    129269 ns/op	     176 B/op	       7 allocs/op
Transition_Denied-8                                        	 2450037	       470 ns/op	     192 B/op	       8 allocs/op
Transition_with_post_transition_functions_run-8         	   10000	    131176 ns/op	     176 B/op	       7 allocs/op
```
-------------
made with love in GO :)
