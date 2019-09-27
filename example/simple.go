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

	// create a new instance of FSM
	fsm, err := go_fsm.NewFsm().
		SetLogger(logger).
		When(
			stateIdle,
			func(eventCtx go_fsm.EventContext, fsmCtx go_fsm.FsmContext) (next go_fsm.State, nextFsmCtx go_fsm.FsmContext, err error) {
				// get current state from eventCtx
				event := go_fsm.EventFromCtx(eventCtx)

				// get current state from fsmContext
				state := go_fsm.StateFromCtx(fsmCtx)

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
				// get current state from eventCtx
				event := go_fsm.EventFromCtx(eventCtx)

				// get current state from fsmContext it must been always but if not then function will return UnknownState constant
				state := go_fsm.StateFromCtx(fsmCtx)

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
