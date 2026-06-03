package errorchan

import (
	"fmt"
	"sync/atomic"
)

// The supported personality modes. Pass them to [SetMode] or [WithMode].
const (
	// ModeDere is the default: sweet, flustered, apologetic, and takes the
	// blame for the failure.
	ModeDere = "dere"
	// ModeTsun is tsundere: annoyed at you, blames your code, and is grudgingly
	// helpful anyway.
	ModeTsun = "tsun"
	// ModeYan is yandere: unsettlingly affectionate about your failures. Do not
	// use in production.
	ModeYan = "yan"
)

// globalMode holds the process-wide default mode as a string. It is read and
// written atomically so [Mode] and [SetMode] are safe for concurrent use.
var globalMode atomic.Value

func init() {
	globalMode.Store(ModeDere)
}

// Mode returns the current global default mode. It is safe to call from
// multiple goroutines.
func Mode() string {
	return globalMode.Load().(string)
}

// SetMode sets the global default mode used when no [WithMode] option is given.
// It returns an error for an unknown mode and leaves the current mode unchanged
// in that case. It is safe to call from multiple goroutines.
func SetMode(mode string) error {
	if !validMode(mode) {
		return fmt.Errorf("errorchan: unknown mode %q", mode)
	}
	globalMode.Store(mode)
	return nil
}
