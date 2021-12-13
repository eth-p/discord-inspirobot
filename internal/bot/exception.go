package bot

import (
	"fmt"
)

// exceptional is a marker type that allows for a throw/catch style of control flow.
//
// This is intended to be used only in very specific circumstances where if-return error handling would otherwise be
// unnecessarily complex and/or verbose.
type exceptional struct {
}

// try runs the closure function and catches any exception instances that were panicked.
// If a non-exception panic occurred, it will be bubbled upwards.
func (c exceptional) try(fn func()) *exception {
	errorChannel := make(chan exception, 1)

	func() {
		// Use defer() and recover() to catch the error.
		defer func() {
			defer close(errorChannel)
			panicErr := recover()
			if panicErr != nil {
				if err, ok := panicErr.(exception); ok {
					errorChannel <- err
					return
				}

				panic(panicErr)
			}
		}()

		// Run the function.
		fn()
	}()

	select {
	case err, ok := <-errorChannel:
		if ok {
			return &err
		}
	}

	return nil
}

// Fatal causes a panic to immediately stop the current function from continuing execution.
// If a nil error is provided, this will panic.
func (c exceptional) Fatal(err error) {
	if err == nil {
		panic("Fatal cannot be called with nil error.")
	}

	panic(exception{err: err})
}

// FatalGuard works similarly to Fatal, but behaves as a no-op when a nil error provided.
func (c exceptional) FatalGuard(err error) {
	if err == nil {
		return
	}

	c.Fatal(err)
}

// exception is an implementation of Go's error type.
// This is thrown by exceptional's Fatal method.
type exception struct {
	err error
}

func (e exception) Unwrap() error {
	return e.err
}

func (e exception) Error() string {
	return fmt.Sprintf("exception: %s", e.err.Error())
}
