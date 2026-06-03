package errorchan

import "fmt"

// Recover restyles a recovered panic into errp. Call it deferred at the top of
// a function that uses a named error return:
//
//	func doThing() (err error) {
//	    defer errorchan.Recover(&err)
//	    // ... code that may panic ...
//	}
//
// If no panic is in flight Recover does nothing and leaves errp untouched. On a
// panic it recovers the value, converts it to an error (a panicked error is
// used directly so its chain and sentinel identity survive), styles it, and
// stores it through errp. A nil errp simply swallows the panic.
func Recover(errp *error, opts ...Option) {
	r := recover()
	if r == nil {
		return
	}

	var err error
	if e, ok := r.(error); ok {
		err = e
	} else {
		err = fmt.Errorf("panic: %v", r)
	}

	if errp != nil {
		*errp = Wrap(err, opts...)
	}
}
