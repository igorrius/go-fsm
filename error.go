package go_fsm

import "errors"

var (
	ErrActionNotFound     = errors.New("action not found")
	ErrCanNotExtractEvent = errors.New("can't extract event from context")
	ErrCanNotExtractState = errors.New("can't extract state from context")
)

func checkErrors(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
