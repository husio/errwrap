package errwrap

import "strings"

type errWrapper struct {
	messages []string
	err      error
}

// Error return strgin representation containing colon separated messages and
// originl error text.
func (e *errWrapper) Error() string {
	if len(e.messages) == 0 {
		return e.err.Error()
	}
	return strings.Join(e.messages, ": ") + ": " + e.err.Error()
}

// Err return wrapped error.
func (e *errWrapper) Err() error {
	return e.err
}

// Err wraps given error with message.
//
// Wrapping without message returns original error.
func Err(err error, msg string) error {
	if msg == "" {
		return err
	}

	if e, ok := err.(*errWrapper); ok {
		// do not modify original error, but reuse all the data
		return &errWrapper{
			messages: append([]string{msg}, e.messages...),
			err:      e.err,
		}
	}

	return &errWrapper{
		messages: []string{msg},
		err:      err,
	}
}

// Is return true if given error is equal to any of given errors. All error
// wrapps are unpacked and the top level error is compared.
func Is(err error, anyof ...error) bool {
	if e, ok := err.(*errWrapper); ok {
		err = e.err
	}

	for _, single := range anyof {
		if e, ok := single.(*errWrapper); ok {
			single = e.err
		}
		if err == single {
			return true
		}
	}

	return false
}
