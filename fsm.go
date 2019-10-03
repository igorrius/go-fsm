package go_fsm

import (
	"context"
	"fmt"
	"sync"
)

type (
	EventContext context.Context
	FsmContext   context.Context

	// Action function is using to describe behavior within transition to the next state and choose what exactly state is needed
	// arguments:
	// eventCtx - context which received with event
	// fsmCtx - global FSM context which transferring between FSM states
	// results:
	// next - the next state where to FSM try to arrive
	// nextFsmCtx - context which will be used in next state
	// err - possible error which can be raised by action handler
	ActionFunc = func(eventCtx EventContext, fsmCtx FsmContext) (next State, nextFsmCtx FsmContext, err error)
)

type Fsm struct {
	logger Logger

	ctx          FsmContext
	state        State
	initialState State

	actionMap             map[State]ActionFunc
	postTransitionFuncMap map[transitionKey][]TransitionFunc

	ctxCancelFunc context.CancelFunc
}

// NewFsm create a new instance of FSM
// if logger was not need should set it to nil
func NewFsm(opts ...Option) *Fsm {
	options := newOptions(opts...)

	fsm := &Fsm{
		actionMap:             map[State]ActionFunc{},
		logger:                options.Logger,
		postTransitionFuncMap: map[transitionKey][]TransitionFunc{},
	}

	return fsm
}

//InitWithState init FSM with initial state
func (fsm *Fsm) InitWithState(state State) (*Fsm, error) {
	if !fsm.isStateExists(state) {
		return nil, fmt.Errorf("invalid initial state [%s]", state)
	}

	fsm.ctx, fsm.ctxCancelFunc = context.WithCancel(context.Background())
	fsm.state, fsm.initialState = state, state
	fsm.logger.Log("Init FSM with state:", state)
	return fsm, nil
}

func (fsm Fsm) isStateExists(state State) bool {
	_, isset := fsm.actionMap[state]
	return isset
}

//CurrentState return FSM current state
func (fsm Fsm) CurrentState() State {
	return fsm.state
}

//When FSM event configuration
func (fsm *Fsm) When(state State, action ActionFunc) *Fsm {
	fsm.actionMap[state] = action
	fsm.logger.Logf("Added an action function for state [%s]", state)
	return fsm
}

//Process event by current state action function
func (fsm *Fsm) ProcessEvent(event Event, eventCtx EventContext) error {
	fsm.logger.Logf("Trying to handle [%s] event", event)
	// check context for nil
	eventCtx = ctxWithEvent(checkAndFixEmptyContext(eventCtx), event)

	// get action function for this state
	f, ok := fsm.actionMap[fsm.state]
	if !ok || f == nil {
		fsm.logger.Logf("Event [%s] for State [%s] processing failed", event, fsm.state)
		return ErrActionNotFound
	}

	// check fsm and event contexts for error before the action call
	if err := checkErrors(fsm.ctx.Err(), eventCtx.Err()); err != nil {
		return err
	}

	// create new context with current state value
	fsmCtx := ctxWithState(fsm.ctx, fsm.state)
	nextState, nextCtx, err := f(eventCtx, fsmCtx)
	if err != nil {
		return err
	}

	// set previous fsm context to next fsm context if nil has been returned by action handler (under the hood magic)
	if nil == nextCtx {
		nextCtx = fsmCtx
	}

	// check fsm, nextFsm and event contexts for error after the action call
	if err := checkErrors(eventCtx.Err(), fsmCtx.Err(), nextCtx.Err()); err != nil {
		return err
	}

	// is next state found?
	if !fsm.isStateExists(nextState) {
		fsm.logger.Logf("State [%s] not found", nextState)
		return ErrActionNotFound
	}

	{
		// create waiting group to sync finish for all async transition functions
		wg := new(sync.WaitGroup)
		// process post state action transition functions [strict to strict]
		if transitionFunctions, ok := fsm.postTransitionFuncMap[newTransitionKey(fsm.state, nextState)]; ok && len(transitionFunctions) > 0 {
			fsm.processTransitionFunctions(wg, nextState, nextCtx, transitionFunctions)
		}

		// process post state action transition functions [strict to any]
		if transitionFunctions, ok := fsm.postTransitionFuncMap[newTransitionKey(fsm.state, "*")]; ok && len(transitionFunctions) > 0 {
			fsm.processTransitionFunctions(wg, nextState, nextCtx, transitionFunctions)
		}

		// process post state action transition functions [any to strict]
		if transitionFunctions, ok := fsm.postTransitionFuncMap[newTransitionKey("*", nextState)]; ok && len(transitionFunctions) > 0 {
			fsm.processTransitionFunctions(wg, nextState, nextCtx, transitionFunctions)
		}

		// process post state action transition functions [any to any]
		if transitionFunctions, ok := fsm.postTransitionFuncMap[newTransitionKey("*", "*")]; ok && len(transitionFunctions) > 0 {
			fsm.processTransitionFunctions(wg, nextState, nextCtx, transitionFunctions)
		}

		// waiting until all transition functions are finished
		wg.Wait()
	}

	// update current state and context
	fsm.state = nextState
	fsm.ctx = nextCtx

	return nil
}

// close main context and stop all events processing (a try to process event always return an error)
func (fsm *Fsm) Close() {
	fsm.ctxCancelFunc()
	fsm.logger.Log("FSM has closed")
}

// reset FSM state to initial state and initial context
func (fsm *Fsm) Reset() error {
	fsm.Close()
	if _, err := fsm.InitWithState(fsm.initialState); err != nil {
		return err
	}

	fsm.logger.Log("FSM has reset")
	return nil
}

//RegisterPostTransitionFunc add a transition function
func (fsm *Fsm) RegisterPostTransitionFunc(fromState, toState State, fn TransitionFunc) *Fsm {
	key := newTransitionKey(fromState, toState)
	fsm.postTransitionFuncMap[key] = append(fsm.postTransitionFuncMap[key], fn)
	return fsm
}

func (fsm *Fsm) processTransitionFunctions(wg *sync.WaitGroup, nextState State, nextCtx FsmContext, transitionFunctions []TransitionFunc) {
	wg.Add(len(transitionFunctions))
	for _, fn := range transitionFunctions {
		go func(from, to State, ctx FsmContext, f TransitionFunc) {
			if err := f(from, to, ctx); err != nil {
				fsm.logger.Logf("Transition function from state [%s] to state [%s] call error [%s]",
					from,
					to,
					err.Error(),
				)
			}
			wg.Done()
		}(fsm.state, nextState, nextCtx, fn)
	}
}
