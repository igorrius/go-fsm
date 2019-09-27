package go_fsm

import "errors"

var (
	ErrActionNotFound = errors.New("action not found")
)

func checkErrors(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
