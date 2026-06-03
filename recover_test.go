package errorchan

import (
	"errors"
	"strings"
	"testing"
)

func panicsWithError(opts ...Option) (err error) {
	defer Recover(&err, opts...)
	panic(errSentinel)
}

func panicsWithString(opts ...Option) (err error) {
	defer Recover(&err, opts...)
	panic("everything is on fire")
}

func returnsNormally() (err error) {
	defer Recover(&err)
	return errors.New("ordinary error")
}

func TestRecoverErrorPanic(t *testing.T) {
	err := panicsWithError(WithMode(ModeTsun), WithSeed(1))
	if err == nil {
		t.Fatal("Recover did not capture the panic")
	}
	if !errors.Is(err, errSentinel) {
		t.Errorf("recovered error lost sentinel identity: %v", err)
	}
}

func TestRecoverStringPanic(t *testing.T) {
	err := panicsWithString(WithSeed(1))
	if err == nil {
		t.Fatal("Recover did not capture the string panic")
	}
	if !strings.Contains(err.Error(), "everything is on fire") {
		t.Errorf("recovered error dropped the panic value: %v", err)
	}
}

func TestRecoverNoPanicLeavesErrUntouched(t *testing.T) {
	err := returnsNormally()
	if err == nil || err.Error() != "ordinary error" {
		t.Errorf("Recover altered the normal return path: %v", err)
	}
}

func TestRecoverNilPointerSwallowsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic escaped Recover(nil): %v", r)
		}
	}()
	func() {
		defer Recover(nil)
		panic("boom")
	}()
}

func TestRecoverOutsidePanicIsNoop(t *testing.T) {
	var err error
	func() { defer Recover(&err) }() // no panic in flight
	if err != nil {
		t.Errorf("Recover set err with no panic: %v", err)
	}
}
