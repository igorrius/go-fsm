package go_fsm

type (
	// transition function is using to add an additional behavior after transition to the next state
	TransitionFunc = func(from, to State, fsmCtx FsmContext) error
)

type transitionKey struct {
	from, to State
}

func newTransitionKey(from State, to State) transitionKey {
	return transitionKey{from: from, to: to}
}
